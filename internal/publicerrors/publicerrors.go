// Package publicerrors provides a way to mark errors as public or private.
// It supplies functions to get the error message of an error, which hides
// possibly sensitive information if the error is private.
//
// The time complexity of checking if an error is public or not is O(n), where n
// is how deep the error value is wrapped.
package publicerrors

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"log/slog"
	"reflect"

	"github.com/lmittmann/tint"
)

var publicErrorValues = map[error]publicError{}

// MarkValuesPublic marks the given errors as public.
func MarkValuesPublic(errors ...error) {
	for _, err := range errors {
		setMapUnique(publicErrorValues, err, publicError{})
	}
}

// MarkValuePublicWithStringer marks the given error as public with a custom
// function for formatting the error.
func MarkValuePublicWithStringer[T error](err T, stringer func(ctx context.Context, outer string) string) {
	setMapUnique(publicErrorValues, error(err), publicError{
		stringer: func(ctx context.Context, err error, outer string) string {
			return stringer(ctx, outer)
		},
	})
}

var publicErrorTypes = map[reflect.Type]publicError{
	reflect.TypeFor[forcedPublicError](): {},
}

type publicError struct {
	stringer func(ctx context.Context, err error, outer string) string
}

// MarkTypePublic marks the given error type as public.
func MarkTypePublic[T error]() {
	setMapUnique(publicErrorTypes, reflect.TypeFor[T](), publicError{})
}

// MarkTypePublicWithStringer marks the given error type as public with a custom
// function for formatting the error.
func MarkTypePublicWithStringer[T error](stringer func(ctx context.Context, err T, outer string) string) {
	setMapUnique(publicErrorTypes, reflect.TypeFor[T](), publicError{
		stringer: func(ctx context.Context, err error, outer string) string {
			return stringer(ctx, err.(T), outer)
		},
	})
}

func setMapUnique[K comparable, V any](m map[K]V, k K, v V) {
	if _, ok := m[k]; ok {
		panic(fmt.Sprintf("duplicate key %#v", k))
	}
	m[k] = v
}

type forcedPublicError struct {
	err error
}

func (e forcedPublicError) Error() string { return e.err.Error() }
func (e forcedPublicError) Unwrap() error { return e.err }

// New creates a new public error with the given text.
func New(text string) error {
	return forcedPublicError{errors.New(text)}
}

// Errorf formats the error message with the given format and arguments.
// The returned error is public.
func Errorf(format string, args ...any) error {
	return forcedPublicError{fmt.Errorf(format, args...)}
}

// ForcePublic forces the given error to be public.
// The returned error will pass all check functions as if it was public.
func ForcePublic(err error) error {
	return forcedPublicError{err}
}

// IsPublic reports if the given error is public.
// This function should rarely be used. Prefer using [String] or [MarshalJSON]
// to get the actual error string instead.
func IsPublic(err error) bool {
	marshaled := MarshalError(context.Background(), err, "")
	return !marshaled.Internal
}

// HiddenMessage is the default message that is shown when a hidden error is
// converted to a string or JSON. The local hiddenMessage variable of a function
// generally overrides this string.
const HiddenMessage = "an internal error occurred :("

// String returns the error message of the given error.
// If the error is public, then the error message is returned. Otherwise, the
// error message is hidden and replaced with [hiddenMessage]. If [hiddenMessage]
// is empty, then the global [HiddenMessage] is used.
func String(ctx context.Context, err error, hiddenMessage string) string {
	marshaled := MarshalError(ctx, err, hiddenMessage)
	return marshaled.Message
}

// MarshaledError is the JSON structure that [MarshalJSON] returns.
// Use this structure to unmarshal the JSON error.
type MarshaledError struct {
	// Message is the error message.
	Message string `json:"message"`
	// Details is the structured error value, if any.
	// It is guaranteed to either be nil or a valid object containing at least 1
	// key.
	Details any `json:"details,omitempty"`
	// Errors is a list of errors, if any.
	// If this is set, then [Details] is ignored.
	Errors []MarshaledError `json:"errors,omitempty"`
	// Internal is true if the error is internal (not public).
	Internal bool `json:"internal,omitempty"`
	// InternalCode is the internal error code, if any.
	InternalCode string `json:"internalCode,omitempty"`
}

// MarshalJSON marshals the given error to JSON.
func MarshalJSON(ctx context.Context, err error, hiddenMessage string) ([]byte, error) {
	marshaled := MarshalError(ctx, err, hiddenMessage)
	return json.Marshal(marshaled)
}

// MarshalError marshals the given error to a [MarshaledError].
// If the error is public, then the error message is used. Otherwise, the error
// message is hidden and replaced with [hiddenMessage]. If err is nil, then the
// error message is [HiddenMessage].
func MarshalError(ctx context.Context, err error, hiddenMessage string) MarshaledError {
	marshaled := MarshaledError{Internal: true}

	var multiError []error
	var multiErrorMode bool

	stack := []error{err}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Check for value.
		public, isPublic := publicErrorValues[curr]
		if !isPublic {
			public, isPublic = publicErrorTypes[reflect.TypeOf(curr)]
		}

		if isPublic {
			marshaled.Internal = false

			if public.stringer != nil {
				marshaled.Message = public.stringer(ctx, curr, marshaled.Message)
			} else {
				if marshaled.Message != "" {
					marshaled.Message = marshaled.Message + ": " + curr.Error()
				} else {
					marshaled.Message = curr.Error()
				}
			}

			if isValidErrorDetails(curr) {
				multiError = append(multiError, curr)
			}
		}

		// If there is no custom stringer, then keep looking for a deeper error.
		// This allows us to use [ForcePublic] while still preserving the
		// stringer of the underlying type.
		switch v := curr.(type) {
		case interface{ Unwrap() error }:
			stack = append(stack, v.Unwrap())
		case interface{ Unwrap() []error }:
			stack = append(stack, v.Unwrap()...)
			multiErrorMode = true
		}
	}

	if multiErrorMode {
		marshaled.Errors = make([]MarshaledError, len(multiError))
		for i, err := range multiError {
			marshaled.Errors[i] = MarshaledError{
				Message: err.Error(),
				Details: err,
			}
		}
	} else if len(multiError) > 0 {
		marshaled.Details = multiError[0]
	}

	if marshaled.Internal {
		marshaled.InternalCode = generateInternalCode(err)
		if marshaled.Message == "" {
			if hiddenMessage == "" {
				marshaled.Message = HiddenMessage
			} else {
				marshaled.Message = hiddenMessage
			}
		}

		slog.ErrorContext(ctx,
			"internal error occured",
			"internal", true,
			"internal_code", marshaled.InternalCode,
			tint.Err(err))
	}

	return marshaled
}

func isValidErrorDetails(err error) bool {
	rv := reflect.Indirect(reflect.ValueOf(err))
	if rv.Kind() != reflect.Struct {
		return false
	}

	rt := rv.Type()

	nfields := rv.NumField()
	for i := 0; i < nfields; i++ {
		fieldv := rv.Field(i)
		fieldt := rt.Field(i)
		if !fieldv.IsZero() && fieldt.IsExported() && fieldt.Tag.Get("json") != "" {
			return true
		}
	}

	return false
}

func generateInternalCode(err error) string {
	h := fnv.New64a()
	h.Write([]byte(err.Error()))
	s := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return s
}

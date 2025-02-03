package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5/middleware"
	"e2clicker.app/internal/publicerrors"
	"e2clicker.app/internal/slogutil"
	"e2clicker.app/internal/userlimit"
	"e2clicker.app/services/api/openapi"
)

func init() {
	publicerrors.MarkTypePublic[*openapi3filter.RequestError]()
	publicerrors.MarkTypePublic[*openapi3filter.ResponseError]()
	publicerrors.MarkTypePublic[*openapi3filter.ParseError]()
	publicerrors.MarkTypePublic[*openapi3filter.ValidationError]()
	publicerrors.MarkTypePublic[*openapi3filter.SecurityRequirementsError]()
	publicerrors.MarkTypePublic[*openapi3.SchemaError]()

	publicerrors.MarkValuesPublic(
		openapi3filter.ErrInvalidRequired,
		openapi3filter.ErrInvalidEmptyValue,
		openapi3filter.ErrAuthenticationServiceMissing,
		ErrNoAcceptableContentType,
	)
}

// ErrNoAcceptableContentType is returned when the request does not have an
// acceptable Content-Type OR Accept header.
var ErrNoAcceptableContentType = errors.New("no acceptable content type")

type errorResponse = struct {
	Body       openapi.Error
	StatusCode int
}

func convertError[T ~errorResponse](ctx context.Context, err error) T {
	return convertErrorWithMessage[T](ctx, err, "")
}

func convertErrorWithMessage[T ~errorResponse](ctx context.Context, err error, hiddenMessage string) T {
	marshaled := publicerrors.MarshalError(ctx, err, hiddenMessage)
	return convertErrorWithMessageFromMarshaled[T](ctx, marshaled)
}

func convertErrorWithMessageFromMarshaled[T ~errorResponse](ctx context.Context, marshaled publicerrors.MarshaledError) T {
	oapiValue := openapi.Error{Message: marshaled.Message}
	if len(marshaled.Errors) > 0 {
		oapiValue.Errors = make([]openapi.Error, len(marshaled.Errors))
		for i, e := range marshaled.Errors {
			resp := convertErrorWithMessageFromMarshaled[errorResponse](ctx, e)
			oapiValue.Errors[i] = resp.Body
		}
	}
	if marshaled.Details != nil {
		oapiValue.Details = &marshaled.Details
	}
	if marshaled.Internal {
		oapiValue.Internal = &marshaled.Internal
		oapiValue.InternalCode = &marshaled.InternalCode
	}

	statusCode := 400
	if marshaled.Internal {
		statusCode = 500
	}

	return T{
		Body:       oapiValue,
		StatusCode: statusCode,
	}
}

func writeError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	resp := convertError[errorResponse](r.Context(), err)
	if !optPtr(resp.Body.Internal) {
		if statusCode != 0 {
			// Error is revealed, so use the suggested status code.
			resp.StatusCode = statusCode
		}

		var rateLimitError *userlimit.LimitExceededError
		if errors.As(err, &rateLimitError) {
			resp.StatusCode = http.StatusTooManyRequests
			resp.Body.Details = anyPtr(rateLimitError)

			if retryAfter, ok := rateLimitError.DelaySeconds(); ok {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
			}
		}
	}

	b, err := json.Marshal(resp.Body)
	if err != nil {
		resp = convertError[errorResponse](r.Context(), err)

		b, err = json.Marshal(resp.Body)
		if err != nil {
			panic(fmt.Errorf("cannot marshal fallback error: %w", err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(b)
}

func recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			v := recover()
			if v == nil {
				return
			}

			slog.ErrorContext(r.Context(),
				"panic recovered while handling request",
				"panic", v,
				slogutil.StackTrace(2))

			if r.Header.Get("Connection") == "Upgrade" || ww.Status() != 0 {
				// Status code has already been written. Don't write
				// anything else.
				ww.Discard()
			}

			writeError(w, r, fmt.Errorf("%v", v), http.StatusInternalServerError)
		}()

		next.ServeHTTP(ww, r)
	})
}

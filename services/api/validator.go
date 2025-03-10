package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"libdb.so/ctxt"
	"e2clicker.app/internal/publicerrors"
	"e2clicker.app/services/user"

	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
)

func init() {
	publicerrors.MarkTypePublic[*routers.RouteError]()
	publicerrors.MarkTypePublic[*validateRequestError]()
	publicerrors.MarkTypePublic[*securityRequirementsError]()

	publicerrors.MarkValuesPublic(
		routers.ErrPathNotFound,
		routers.ErrMethodNotAllowed,
	)
}

// requestContextOverrider allows functions within [validateRequest] to override
// the request's context using the parent context that it was given by providing
// a pointer value that functions can mutate from within.
//
// This is a workaround for the fact that [openapi3filter.ValidateRequest] does
// not allow functions to override the context.
type requestContextOverrider struct {
	old context.Context
	new context.Context
}

func addToRequestContext[T any](ctx context.Context, v T) {
	co, ok := ctxt.From[*requestContextOverrider](ctx)
	if !ok {
		panic("ctx is not used within the RequestValidator middleware")
	}
	co.old = co.new
	co.new = ctxt.With[T](ctx, v)
}

func newRequestValidator(spec *openapi3.T, options openapi3filter.Options) func(next http.Handler) http.Handler {
	router, err := legacyrouter.NewRouter(spec)
	if err != nil {
		panic(fmt.Errorf("cannot build openapi3 router from spec: %w", err))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			co := &requestContextOverrider{
				old: r.Context(),
				new: r.Context(),
			}
			r = r.WithContext(ctxt.With(r.Context(), co))

			statusCode, err := validateRequest(r, router, options)
			if err != nil {
				writeError(w, r, err, statusCode)
				return
			}

			r = r.WithContext(co.new)
			next.ServeHTTP(w, r)
		})
	}
}

type validateRequestError struct {
	Reason      string                `json:"reason"`
	Err         string                `json:"error"`
	Parameter   *openapi3.Parameter   `json:"parameter,omitempty"`
	RequestBody *openapi3.RequestBody `json:"requestBody,omitempty"`
}

func (e *validateRequestError) Error() string {
	if e.Err != "" {
		return e.Err
	}
	return e.Reason
}

type securityRequirementsError struct {
	SecurityRequirements []openapi3.SecurityRequirement `json:"securityRequirements"`
	Errors               []string                       `json:"errors"`
}

func (e *securityRequirementsError) Error() string {
	s := "security requirements not satisfied"
	if len(e.Errors) > 0 {
		s += ": " + strings.Join(e.Errors, ", ")
	}
	return s
}

// validateRequest is called from the middleware above and actually does the work
// of validating a request.
func validateRequest(r *http.Request, router routers.Router, options openapi3filter.Options) (int, error) {
	route, pathParams, err := router.FindRoute(r)
	if err != nil {
		return http.StatusNotFound, err
	}

	if !routeBodyContainsApplication(route) {
		// Disable body validation
		options.ExcludeRequestBody = true
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
		Options:    &options,
	}

	if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
		switch e := err.(type) {
		case *openapi3filter.RequestError:
			err := &validateRequestError{
				Parameter:   e.Parameter,
				RequestBody: e.RequestBody,
				Reason:      e.Reason,
			}
			if e.Err != nil {
				err.Err = e.Err.Error()
				err.Err, _, _ = strings.Cut(err.Err, "\n")
			}
			return http.StatusBadRequest, err
		case *openapi3filter.SecurityRequirementsError:
			errorStrings := make([]string, len(e.Errors))
			for i, err := range e.Errors {
				errorStrings[i] = err.Error()
			}
			err := &securityRequirementsError{
				SecurityRequirements: e.SecurityRequirements,
				Errors:               errorStrings,
			}
			status := http.StatusForbidden
			if len(e.Errors) == 1 && errors.Is(e.Errors[0], user.ErrInvalidSession) {
				status = http.StatusUnauthorized
			}
			return status, err
		default:
			return http.StatusBadRequest, publicerrors.ForcePublic(err)
		}
	}

	return http.StatusOK, nil
}

func routeBodyContainsApplication(route *routers.Route) bool {
	if false ||
		route.Operation == nil ||
		route.Operation.RequestBody == nil ||
		route.Operation.RequestBody.Value == nil ||
		route.Operation.RequestBody.Value.Content == nil {
		return false
	}

	for t := range route.Operation.RequestBody.Value.Content {
		if strings.HasPrefix(t, "application/") {
			return true
		}
	}

	return false
}

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"libdb.so/e2clicker/internal/publicerrors"

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

func newRequestValidator(spec *openapi3.T, options *openapi3filter.Options) func(next http.Handler) http.Handler {
	router, err := legacyrouter.NewRouter(spec)
	if err != nil {
		panic(fmt.Errorf("cannot build openapi3 router from spec: %w", err))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			statusCode, err := validateRequest(r, router, options)
			if err != nil {
				writeError(w, r, err, statusCode)
				return
			}
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
func validateRequest(r *http.Request, router routers.Router, options *openapi3filter.Options) (int, error) {
	route, pathParams, err := router.FindRoute(r)
	if err != nil {
		return http.StatusNotFound, err
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
		Options:    options,
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
			return http.StatusForbidden, err
		default:
			return http.StatusBadRequest, publicerrors.ForcePublic(err)
		}
	}

	return http.StatusOK, nil
}

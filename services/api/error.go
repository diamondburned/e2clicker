package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5/middleware"
	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/e2clicker/services/api/openapi"
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
	)
}

type errorResponse = struct {
	Body       openapi.Error
	StatusCode int
}

func convertError[T ~errorResponse](ctx context.Context, err error) T {
	return convertErrorWithMessage[T](ctx, err, "")
}

func convertErrorWithMessage[T ~errorResponse](ctx context.Context, err error, hiddenMessage string) T {
	marshaled := publicerrors.MarshalError(ctx, err, hiddenMessage)
	oapiValue := openapi.Error{Message: marshaled.Message}
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
	errResponse := convertError[errorResponse](r.Context(), err)
	if statusCode != 0 && !optPtr(errResponse.Body.Internal) {
		// Error is revealed, so use the suggested status code.
		errResponse.StatusCode = statusCode
	}

	b, err := json.Marshal(errResponse.Body)
	if err != nil {
		errResponse = convertError[errorResponse](r.Context(), err)

		b, err = json.Marshal(errResponse.Body)
		if err != nil {
			panic(fmt.Errorf("cannot marshal fallback error: %w", err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResponse.StatusCode)
	w.Write(b)
}

func recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			if v := recover(); v != nil {
				if r.Header.Get("Connection") == "Upgrade" || ww.Status() != 0 {
					// Status code has already been written. Don't write
					// anything else.
					ww.Discard()
				}
				writeError(w, r, fmt.Errorf("%v", v), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(ww, r)
	})
}

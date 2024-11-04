package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/fxhooking"
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/hserve"

	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	e2clickermodule "libdb.so/e2clicker/nix/modules/e2clicker"
)

// Server provides an HTTP server that serves a [Handler].
type Server struct {
}

// ServerInputs is a set of dependencies required by the [Server].
type ServerInputs struct {
	fx.In

	Handler       openapi.StrictServerInterface
	Authenticator openapi3filter.AuthenticationFunc
}

// NewServer creates a new HTTP server.
func NewServer(
	lx fx.Lifecycle,
	inputs ServerInputs,
	config e2clickermodule.API,
	logger *slog.Logger,
) (*Server, error) {
	logger = logger.With("addr", config.ListenAddress)

	swaggerAPI, err := openapi.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("cannot get swagger schema: %w", err)
	}

	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(swaggerAPI, &nethttpmiddleware.Options{
		ErrorHandler: writeValidationError,
		Options: openapi3filter.Options{
			MultiError:         false,
			AuthenticationFunc: inputs.Authenticator,
		},
	})

	router := chi.NewRouter()
	router.Use(logRequest(logger))
	router.Use(recovererMiddleware)

	openapi.HandlerWithOptions(
		openapi.NewStrictHandlerWithOptions(
			inputs.Handler,
			[]strictnethttp.StrictHTTPMiddlewareFunc{
				func(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
					return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
						h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							response, err = f(ctx, w, r, request)
						})
						validator(h).ServeHTTP(w, r)
						return
					}
				},
			},
			openapi.StrictHTTPServerOptions{
				RequestErrorHandlerFunc:  errorWriterForCode(http.StatusBadRequest),
				ResponseErrorHandlerFunc: errorWriterForCode(http.StatusInternalServerError),
			},
		),
		openapi.ChiServerOptions{
			BaseURL:          "/api",
			BaseRouter:       router,
			ErrorHandlerFunc: errorWriterForCode(http.StatusBadRequest),
		},
	)

	lx.Append(fxhooking.WrapRun(func(ctx context.Context) error {
		logger.Info("listening to HTTP server")
		defer logger.Warn("HTTP server stopped")

		if err := hserve.ListenAndServe(ctx, config.ListenAddress, router); err != nil {
			return fmt.Errorf("HTTP server error: %w", err)
		}

		return nil
	}))

	return &Server{}, nil
}

func respond200(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

func logRequest(slog *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.DebugContext(r.Context(),
				"received API request",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.Query().Encode(),
				"headers", r.Header)

			next.ServeHTTP(w, r)
		})
	}
}

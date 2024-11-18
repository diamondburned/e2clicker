package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/fxhooking"
	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/hserve"

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

	requestValidatorMiddleware := newRequestValidator(swaggerAPI, openapi3filter.Options{
		ExcludeResponseBody: true,
		AuthenticationFunc:  inputs.Authenticator,
	})

	handler := openapi.HandlerWithOptions(
		openapi.NewStrictHandlerWithOptions(
			inputs.Handler, nil,
			openapi.StrictHTTPServerOptions{
				RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					err = publicerrors.ForcePublic(err) // only validation errors
					writeError(w, r, err, http.StatusBadRequest)
				},
				ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					writeError(w, r, err, 0)
				},
			},
		),
		openapi.StdHTTPServerOptions{
			BaseURL: "/api",
			ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				writeError(w, r, err, http.StatusBadRequest)
			},
		},
	)

	handler = logRequest(logger)(handler)
	handler = requestValidatorMiddleware(handler)
	handler = recovererMiddleware(handler)

	mux := http.NewServeMux()
	mountDocs(mux, swaggerAPI)
	mux.Handle("/api/", handler)

	lx.Append(fxhooking.WrapRun(func(ctx context.Context) error {
		logger.Info("listening to HTTP server")
		defer logger.Warn("HTTP server stopped")

		if err := hserve.ListenAndServe(ctx, config.ListenAddress, mux); err != nil {
			return fmt.Errorf("HTTP server error: %w", err)
		}

		return nil
	}))

	return &Server{}, nil
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

package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/fxhooking"
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/hserve"

	e2clickermodule "libdb.so/e2clicker/nix/modules/e2clicker"
)

// Server provides an HTTP server that serves a [Handler].
type Server struct {
}

// NewServer creates a new HTTP server.
func NewServer(
	lx fx.Lifecycle,
	handler *OpenAPIHandler,
	logger *slog.Logger,
	config e2clickermodule.API,
) (*Server, error) {
	logger = logger.With("addr", config.ListenAddress)

	router := chi.NewRouter()
	router.Use(logRequest(logger))

	openapi.HandlerFromMuxWithBaseURL(
		openapi.NewStrictHandler(handler, nil),
		router, "/api")

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

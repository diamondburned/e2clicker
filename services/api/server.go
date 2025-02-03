package api

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"go.uber.org/fx"
	"e2clicker.app/internal/fxhooking"
	"e2clicker.app/services/api/openapi"
	"libdb.so/hserve"

	e2clickermodule "e2clicker.app/nix/modules/e2clicker"
)

// Server provides an HTTP server that serves a [Handler].
type Server struct {
}

// ServerInputs is a set of dependencies required by the [Server].
type ServerInputs struct {
	fx.In

	Handler       openapi.ServerInterface
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
		inputs.Handler,
		openapi.StdHTTPServerOptions{
			BaseURL: "/api",
			ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				writeError(w, r, err, http.StatusBadRequest)
			},
		},
	)

	handler = applyMiddleware(handler, []middlewareFunc{
		logRequest(logger), // runs first
		recovererMiddleware,
		requestValidatorMiddleware,
	})

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

type middlewareFunc = func(http.Handler) http.Handler

func applyMiddleware(next http.Handler, mws []middlewareFunc) http.Handler {
	for _, mw := range slices.Backward(mws) {
		next = mw(next)
	}
	return next
}

func logRequest(slog *slog.Logger) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggedHeader := cleanupLoggedHeaders(r.Header)

			slog.DebugContext(r.Context(),
				"received API request",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.Query().Encode(),
				"headers", loggedHeader)

			start := time.Now()
			next.ServeHTTP(w, r)
			since := time.Since(start)

			slog.DebugContext(r.Context(),
				"finished API request",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.Query().Encode(),
				"headers", loggedHeader,
				"took", since)
		})
	}
}

func cleanupLoggedHeaders(h http.Header) http.Header {
	h = maps.Clone(h)

	// Censor the Authorization token.
	if h.Get("Authorization") != "" {
		h.Set("Authorization", "***")
	}

	// Delete potentially identifying headers.
	delete(h, "User-Agent")
	for k := range h {
		if strings.HasPrefix(k, "Sec-") {
			delete(h, k)
		}
	}

	return h
}

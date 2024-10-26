package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/hserve"
)

// ServerConfig defines the configuration for [Server].
type ServerConfig struct {
	ListenAddress string `json:"listenAddress"`
}

// Server provides an HTTP server that serves a [Handler].
type Server struct {
	config *ServerConfig `do:""`
	logger *slog.Logger  `do:""`

	router *chi.Mux

	ctx  context.Context
	stop func(error)
	done chan struct{}
}

var (
	_ do.ShutdownerWithContextAndError = (*Server)(nil)
	_ do.HealthcheckerWithContext      = (*Server)(nil)
)

func newServer(i do.Injector) (*Server, error) {
	ctx := do.MustInvoke[context.Context](i)
	s := do.MustInvokeStruct[Server](i)

	s.router = chi.NewRouter()
	s.router.Use(logRequest(s.logger))

	oapiHandler := do.MustInvokeStruct[oapiHandler](i)
	openapi.HandlerFromMux(
		openapi.NewStrictHandler(oapiHandler, nil),
		s.router)

	s.done = make(chan struct{})
	s.ctx, s.stop = context.WithCancelCause(ctx)

	go func() {
		defer close(s.done)

		s.logger.Info(
			"listening to HTTP server",
			"addr", s.config.ListenAddress)

		if err := hserve.ListenAndServe(s.ctx, s.config.ListenAddress, s.router); err != nil {
			s.logger.Error("http server error", "err", err)
			s.stop(err)
		}
	}()

	return s, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.stop(nil)
	select {
	case <-s.done:
	case <-ctx.Done():
		return ctx.Err()
	}

	cause := context.Cause(s.ctx)
	if cause != nil && cause != s.ctx.Err() {
		return cause
	}

	return nil
}

func (s *Server) HealthCheck(ctx context.Context) error {
	select {
	case <-s.done:
		return context.Cause(s.ctx)
	default:
		return nil
	}
}

func respond200(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func logRequest(logger *slog.Logger) func(next http.Handler) http.Handler {
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

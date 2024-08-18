package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/lmittmann/tint"
	"github.com/plaid/go-envvar/envvar"
	"github.com/spf13/pflag"
	"libdb.so/e2clicker/cmd/e2clicker-backend/cfgtypes"
	"libdb.so/e2clicker/services/storage/postgresql"
	"libdb.so/hserve"
)

var (
	logFormat = cfgtypes.NewStringEnum("color", "json", "text")
	verbosity = 0
)

var env struct {
	DatabaseURI string `envvar:"DATABASE_URI"`
	HTTPAddress string `envvar:"HTTP_ADDRESS"`
}

var logLevel = slog.LevelWarn

func init() {
	log.SetFlags(0)

	pflag.VarP(logFormat, "log-format", "l", "Log format (color, json, text)")
	pflag.CountVarP(&verbosity, "verbose", "v", "Increase verbosity (default: warn)")
}

func main() {
	pflag.Parse()
	if err := envvar.Parse(&env); err != nil {
		log.Fatalln(err)
	}
	logLevel -= slog.Level(verbosity) * 4 // every 4 level is a constant

	var slogHandler slog.Handler
	switch logFormat.Value {
	case "color":
		slogHandler = tint.NewHandler(os.Stderr, &tint.Options{Level: logLevel})
	case "json":
		slogHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	case "text":
		slogHandler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	default:
		panic("unknown log format")
	}

	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		cancel()
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	dbURI, err := url.Parse(env.DatabaseURI)
	if err != nil {
		return fmt.Errorf("invalid database URI: %w", err)
	}

	switch dbURI.Scheme {
	case "postgres", "postgresql":
		_, err := postgresql.Connect(ctx, env.DatabaseURI)
		if err != nil {
			return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
	default:
		return fmt.Errorf("unsupported database scheme %q in URI", dbURI.Scheme)
	}

	mux := chi.NewMux()
	if logLevel <= slog.LevelDebug {
		mux.Use(logRequest)
	}

	mux.Get("/", respond200)

	slog.Info(
		"listening to HTTP server",
		"addr", env.HTTPAddress)

	if err := hserve.ListenAndServe(ctx, env.HTTPAddress, mux); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func respond200(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func logRequest(next http.Handler) http.Handler {
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

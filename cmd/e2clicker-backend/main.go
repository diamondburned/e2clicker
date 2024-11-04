package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"libdb.so/e2clicker/services/api"
	"libdb.so/e2clicker/services/notification"
	"libdb.so/e2clicker/services/storage"
	"libdb.so/e2clicker/services/user"

	e2clickermodule "libdb.so/e2clicker/nix/modules/e2clicker"
)

var (
	configFile       = "config.json"
	explainRootScope = false
	explainService   = false
)

var logLevel = slog.LevelWarn

func init() {
	log.SetFlags(0)

	pflag.StringVarP(&configFile, "config", "c", configFile, "Configuration file")
	pflag.BoolVar(&explainRootScope, "explain-root-scope", explainRootScope, "Explain the root scope")
	pflag.BoolVar(&explainService, "explain-service", explainService, "Explain the service")
}

func main() {
	pflag.Parse()

	var cfg e2clickermodule.BackendConfig
	if err := readJSONFile(configFile, &cfg); err != nil {
		log.Fatalln("cannot read config file:", err)
	}

	logLevel := slog.LevelInfo
	if cfg.Debug {
		logLevel = slog.LevelDebug
	}

	var slogHandler slog.Handler
	switch cfg.LogFormat {
	case e2clickermodule.LogFormatColor:
		slogHandler = tint.NewHandler(os.Stderr, &tint.Options{Level: logLevel})
	case e2clickermodule.LogFormatJSON:
		slogHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	case e2clickermodule.LogFormatText:
		slogHandler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	default:
		panic("unknown log format")
	}

	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	fx.New(
		api.Module,
		user.Module,
		storage.Module,
		notification.Module,
		fx.Supply(slog.Default()),
		fx.Supply(cfg.API),
		fx.Supply(cfg.PostgreSQL),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			l := &fxevent.SlogLogger{Logger: logger}
			l.UseLogLevel(slog.LevelDebug)
			return l
		}),
		fx.Invoke(func(*api.Server) {
			slog.Info("API server started successfully")
		}),
	).Run()
}

func readJSONFile(path string, dst any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	if err := json.Unmarshal(b, dst); err != nil {
		return fmt.Errorf("cannot unmarshal config: %w", err)
	}

	return nil
}

func indent(s string) string {
	l := strings.Split(s, "\n")
	for i, line := range l {
		l[i] = "\t" + line
	}
	return strings.Join(l, "\n")
}

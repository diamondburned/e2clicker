package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/samber/do/v2"
	"github.com/spf13/pflag"
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

	if err := run(context.Background(), cfg); err != nil {
		slog.Error(
			"program failed",
			tint.Err(err))
		os.Exit(1)
	}
}

var Package = do.Package(
	api.Package,
	user.Package,
	storage.Package,
	notification.Package,
)

func run(ctx context.Context, cfg e2clickermodule.BackendConfig) error {
	root := do.NewWithOpts(&do.InjectorOpts{
		Logf: func(s string, args ...interface{}) {
			s = strings.TrimPrefix(s, "DI: ")
			slog.DebugContext(ctx, fmt.Sprintf(s, args...), "module", "do")
		},
	}, Package)

	do.ProvideValue(root, ctx)
	do.ProvideValue(root, slog.Default())
	do.ProvideValue(root, cfg.API)
	do.ProvideValue(root, cfg.PostgreSQL)

	if explainRootScope {
		explanation := do.ExplainInjector(root)
		fmt.Println(explanation.String())
		return nil
	}

	_, err := do.Invoke[*api.Server](root)
	if err != nil {
		return err
	}

	if explainService {
		explanation, _ := do.ExplainService[*api.Server](root)
		fmt.Println(explanation.String())
	}

	_, errs := root.ShutdownOnSignalsWithContext(ctx, os.Interrupt)
	if errs != nil && errs.Len() > 0 {
		return errs
	}

	return nil
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

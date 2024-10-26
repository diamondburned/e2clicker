package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/samber/do/v2"
	"github.com/spf13/pflag"
	"libdb.so/e2clicker/cmd/e2clicker-backend/cfgtypes"
	"libdb.so/e2clicker/services/api"
	"libdb.so/e2clicker/services/notification"
	"libdb.so/e2clicker/services/storage"
	"libdb.so/e2clicker/services/user"
)

var (
	configFile       = "config.json"
	logFormat        = cfgtypes.NewStringEnum("color", "json", "text")
	verbosity        = 0
	explainRootScope = false
	explainService   = false
)

var logLevel = slog.LevelWarn

func init() {
	log.SetFlags(0)

	pflag.StringVarP(&configFile, "config", "c", configFile, "Configuration file")
	pflag.VarP(logFormat, "log-format", "l", "Log format (color, json, text)")
	pflag.CountVarP(&verbosity, "verbose", "v", "Increase verbosity (default: warn)")
	pflag.BoolVar(&explainRootScope, "explain-root-scope", explainRootScope, "Explain the root scope")
	pflag.BoolVar(&explainService, "explain-service", explainService, "Explain the service")
}

func main() {
	pflag.Parse()
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

	if err := run(context.Background()); err != nil {
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

func run(ctx context.Context) error {
	root := do.NewWithOpts(&do.InjectorOpts{
		Logf: func(s string, args ...interface{}) {
			s = strings.TrimPrefix(s, "DI: ")
			m := fmt.Sprintf(s, args...)
			slog.DebugContext(ctx, m, "module", "do")
		},
	}, Package)

	do.ProvideValue(root, ctx)
	do.ProvideValue(root, slog.Default())

	if err := parseFile(root, configFile); err != nil {
		return err
	}

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

func indent(s string) string {
	l := strings.Split(s, "\n")
	for i, line := range l {
		l[i] = "\t" + line
	}
	return strings.Join(l, "\n")
}

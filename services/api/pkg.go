package api

import (
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module("api",
	fx.Decorate(func(slog *slog.Logger) *slog.Logger {
		return slog.With("module", "api")
	}),
	fx.Provide(
		NewOpenAPIHandler,
		NewServer,
	),
)

package user

import (
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Decorate(func(slog *slog.Logger) *slog.Logger {
		return slog.With("module", "user")
	}),
	fx.Provide(
		NewUserService,
	),
)

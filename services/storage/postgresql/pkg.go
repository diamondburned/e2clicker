package postgresql

import (
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module("postgresql",
	fx.Decorate(func(slog *slog.Logger) *slog.Logger {
		return slog.With("module", "postgresql")
	}),
	fx.Provide(
		NewStorage,
		(*Storage).userStorage,
		(*Storage).userAvatarStorage,
		(*Storage).userSessionStorage,
		(*Storage).notificationUserStorage,
		(*Storage).dosageStorage,
	),
)

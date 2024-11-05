package dosage

import (
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module("dosage",
	fx.Decorate(func(slog *slog.Logger) *slog.Logger {
		return slog.With("module", "dosage")
	}),
)

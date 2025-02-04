// Package notification provides a way to send notifications to users.
package notification

import (
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module("notification",
	fx.Decorate(func(slog *slog.Logger) *slog.Logger {
		return slog.With("module", "notification")
	}),
	fx.Provide(
		NewNotificationService,
		NewUserNotificationService,
		NewGotifyService,
		NewPushoverService,
		NewWebPushSevice,
		NewEmailService,
	),
)

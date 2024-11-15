package notification

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/validating"
)

// NotificationConfigs contains all the configurations for a notification.
type NotificationConfigs struct {
	Gotify   []GotifyNotificationConfig   `json:"gotify,omitempty"`
	Pushover []PushoverNotificationConfig `json:"pushover,omitempty"`
	WebPush  []WebPushNotificationConfig  `json:"webPush,omitempty"`
}

// IsEmpty returns true if the notification configs are empty.
func (c NotificationConfigs) IsEmpty() bool {
	return len(c.Gotify) == 0 && len(c.Pushover) == 0 && len(c.WebPush) == 0
}

// NotificationService is a collection of NotificationServices.
// It implements the [Notifier] interface.
type NotificationService struct {
	services     NotificationServiceConfig
	servicesAttr slog.Attr
	logger       *slog.Logger
}

// NotificationServiceConfig is the configuration for the notification service.
type NotificationServiceConfig struct {
	fx.In
	Gotify   *GotifyService   `optional:"true"`
	Pushover *PushoverService `optional:"true"`
	WebPush  *WebPushService  `optional:"true"`
}

// NewNotificationService creates a new notification service.
func NewNotificationService(s NotificationServiceConfig, logger *slog.Logger) *NotificationService {
	return &NotificationService{
		services: s,
		logger:   logger,
		servicesAttr: slog.Group(
			"services",
			"gotify", s.Gotify != nil,
			"pushover", s.Pushover != nil,
			"webPush", s.WebPush != nil,
		),
	}
}

func (m *NotificationService) Notify(ctx context.Context, n Notification, c NotificationConfigs) error {
	return errors.Join(slices.Concat(
		callNotify(ctx, "gotify", n, c.Gotify, m.services.Gotify),
		callNotify(ctx, "pushover", n, c.Pushover, m.services.Pushover),
		callNotify(ctx, "webPush", n, c.WebPush, m.services.WebPush),
	)...)
}

func callNotify[
	ConfigT any,
	NotifierT interface {
		Notify(context.Context, Notification, ConfigT) error
	},
](
	ctx context.Context,
	name string,
	notification Notification,
	configs []ConfigT,
	notifier *NotifierT,
) (errs []error) {
	if notifier == nil {
		errs = append(errs, UnknownServiceError{name})
		return
	}
	for _, c := range configs {
		if err := validating.ShouldValidate(ctx, c); err != nil {
			errs = append(errs, ConfigError{"gotify", err})
			continue
		}
		if err := (*notifier).Notify(ctx, notification, c); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
	}
	return
}

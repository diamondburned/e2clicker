package notification

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"e2clicker.app/internal/validating"
	"e2clicker.app/services/notification/openapi"
	"go.uber.org/fx"
)

// NotificationConfigs contains all the configurations for a notification.
type NotificationConfigs struct {
	Gotify   []GotifyNotificationConfig   `json:"gotify,omitempty"`
	Pushover []PushoverNotificationConfig `json:"pushover,omitempty"`
	WebPush  []openapi.PushSubscription   `json:"webPush,omitempty"`
	Email    []EmailNotificationConfig    `json:"email,omitempty"`
}

// NotificationMethodSupports lists the supported notification services.
type NotificationMethodSupports struct {
	Gotify   bool `json:"gotify"`
	Pushover bool `json:"pushover"`
	WebPush  bool `json:"webPush"`
	Email    bool `json:"email"`
}

// IsEmpty returns true if the notification configs are empty.
func (c NotificationConfigs) IsEmpty() bool {
	return len(c.Gotify) == 0 && len(c.Pushover) == 0 && len(c.WebPush) == 0 && len(c.Email) == 0
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
	Email    *EmailService    `optional:"true"`
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
			"email", s.Email != nil,
		),
	}
}

// Notify sends a notification to all the services.
func (m *NotificationService) Notify(ctx context.Context, n Notification, c NotificationConfigs) error {
	return errors.Join(slices.Concat(
		callNotify(ctx, "gotify", n, c.Gotify, m.services.Gotify),
		callNotify(ctx, "pushover", n, c.Pushover, m.services.Pushover),
		callNotify(ctx, "webPush", n, c.WebPush, m.services.WebPush),
		callNotify(ctx, "email", n, c.Email, m.services.Email),
	)...)
}

// Supports returns the supported notification services.
func (m *NotificationService) Supports() NotificationMethodSupports {
	return NotificationMethodSupports{
		Gotify:   m.services.Gotify != nil,
		Pushover: m.services.Pushover != nil,
		WebPush:  m.services.WebPush != nil,
		Email:    m.services.Email != nil,
	}
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
		// errs = append(errs, UnknownServiceError{name})
		return
	}
	for _, c := range configs {
		if err := validating.ShouldValidate(ctx, c); err != nil {
			errs = append(errs, ConfigError{name, err})
			continue
		}
		if err := (*notifier).Notify(ctx, notification, c); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
	}
	return
}

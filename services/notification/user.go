package notification

import (
	"context"

	"go.uber.org/fx"
	"libdb.so/e2clicker/services/user"
)

// UserPreferences is the preferences of a user.
// It is JSON-serializable.
type UserPreferences struct {
	NotificationConfig  NotificationConfigJSON                   `json:"notificationConfig"`
	CustomNotifications map[NotificationMessageType]Notification `json:"customNotification,omitempty"`
}

type UserStorage interface {
	// UserPreferences returns the preferences of a user.
	// If null, the user will not receive any notifications.
	UserPreferences(ctx context.Context, userSecret user.Secret) (*UserPreferences, error)
	// SetUserPreferences sets the preferences of a user.
	// Null is an acceptable value and means the user will not receive any
	// notifications.
	SetUserPreferences(ctx context.Context, userSecret user.Secret, prefs *UserPreferences) error
}

// UserNotificationService is a service that sends notifications to users.
type UserNotificationService struct {
	notifier Notifier
	storage  UserStorage
}

// UserNotificationServiceConfig is the configuration for the user notification
// service.
type UserNotificationServiceConfig struct {
	fx.In

	Notifier
	UserStorage
}

// NewUserNotificationService creates a new user notification service.
func NewUserNotificationService(c UserNotificationServiceConfig) *UserNotificationService {
	return &UserNotificationService{
		notifier: c.Notifier,
		storage:  c.UserStorage,
	}
}

// NotifyUser sends a notification to a user.
func (u *UserNotificationService) NotifyUser(ctx context.Context, id user.Secret, t NotificationMessageType) error {
	prefs, err := u.storage.UserPreferences(ctx, id)
	if err != nil {
		return err
	}

	if prefs == nil || prefs.NotificationConfig.NotificationConfig == nil {
		return nil
	}

	var n *Notification
	if custom, ok := prefs.CustomNotifications[t]; ok {
		n = &custom
	} else {
		n, err = LoadNotification(ctx, t)
		if err != nil {
			return err
		}
	}

	return u.notifier.Notify(ctx, n, prefs.NotificationConfig)
}

// UserPreferences returns the preferences of a user.
// If null, the user will not receive any notifications.
func (u *UserNotificationService) UserPreferences(ctx context.Context, secret user.Secret) (*UserPreferences, error) {
	p, err := u.storage.UserPreferences(ctx, secret)
	if err != nil || p != nil {
		return p, err
	}
	return &UserPreferences{}, nil
}

// SetUserPreferences sets the preferences of a user.
// Null is an acceptable value and means the user will not receive any
// notifications.
func (u *UserNotificationService) SetUserPreferences(ctx context.Context, secret user.Secret, prefs *UserPreferences) error {
	return u.storage.SetUserPreferences(ctx, secret, prefs)
}

// DeleteUserPreferences calls [SetUserPreferences] with a nil [UserPreferences].
func (u *UserNotificationService) DeleteUserPreferences(ctx context.Context, secret user.Secret) error {
	return u.storage.SetUserPreferences(ctx, secret, nil)
}

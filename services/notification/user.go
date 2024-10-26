package notification

import (
	"context"

	"libdb.so/e2clicker/services/user"
)

// UserPreferences is the preferences of a user.
type UserPreferences struct {
	NotificationConfig  NotificationConfig                       `json:"notification_config,omitempty"`
	CustomNotifications map[NotificationMessageType]Notification `json:"custom_notification,omitempty"`
}

type UserStorage interface {
	// UserPreferences returns the preferences of a user.
	// If null, the user will not receive any notifications.
	UserPreferences(ctx context.Context, id user.UserID) (*UserPreferences, error)
	// SetUserPreferences sets the preferences of a user.
	SetUserPreferences(ctx context.Context, id user.UserID, prefs *UserPreferences) error
	// DeleteUserPreferences deletes the preferences of a user.
	DeleteUserPreferences(ctx context.Context, id user.UserID) error
}

type UserNotificationService struct {
	notifier Notifier    `do:""`
	storage  UserStorage `do:""`
}

// NotifyUser sends a notification to a user.
func (u *UserNotificationService) NotifyUser(ctx context.Context, id user.UserID, t NotificationMessageType) error {
	prefs, err := u.storage.UserPreferences(ctx, id)
	if err != nil {
		return err
	}

	if prefs == nil || prefs.NotificationConfig == nil {
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

// UserPreferences gets the preferences of a user.
// The function never returns nil.
func (u *UserNotificationService) UserPreferences(ctx context.Context, id user.UserID) (*UserPreferences, error) {
	p, err := u.storage.UserPreferences(ctx, id)
	if err != nil || p != nil {
		return p, err
	}
	return &UserPreferences{}, nil
}

func (u *UserNotificationService) SetUserPreferences(ctx context.Context, id user.UserID, prefs *UserPreferences) error {
	if prefs == nil {
		prefs = &UserPreferences{}
	}
	return u.storage.SetUserPreferences(ctx, id, prefs)
}

func (u *UserNotificationService) DeleteUserPreferences(ctx context.Context, id user.UserID) error {
	return u.storage.DeleteUserPreferences(ctx, id)
}

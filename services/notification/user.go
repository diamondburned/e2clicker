package notification

import (
	"context"

	"libdb.so/e2clicker/services/user"
)

type UserNotificationStorage interface {
	// SetUserNotificationService sets the notification service for a user.
	// If null, the user will not receive any notifications.
	SetUserNotificationService(ctx context.Context, id user.UserID, config *NotificationConfig) error
	// SetUserCustomNotification sets a custom notification for a user.
	// If null, then the default notification will be used.
	SetUserCustomNotification(ctx context.Context, id user.UserID, n *Notification) error
}

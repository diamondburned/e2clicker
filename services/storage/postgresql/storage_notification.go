package postgresql

import (
	"context"

	"libdb.so/e2clicker/internal/sqlc"
	"libdb.so/e2clicker/internal/sqlc/postgresqlc"
	"libdb.so/e2clicker/services/notification"
	"libdb.so/e2clicker/services/user"
)

func (s *Storage) notificationUserStorage() notification.UserNotificationStorage {
	return (*notificationUserStorage)(s)
}

type notificationUserStorage Storage

func (s *notificationUserStorage) UserPreferences(ctx context.Context, userSecret user.Secret) (notification.UserPreferences, error) {
	p, err := s.q.UserNotificationPreferences(ctx, sqlc.XID(userSecret))
	if err != nil || p == nil {
		return notification.UserPreferences{}, err
	}
	return *p, nil
}

func (s *notificationUserStorage) SetUserPreferences(ctx context.Context, userSecret user.Secret, prefs notification.UserPreferences) error {
	return s.q.SetUserNotificationPreferences(ctx, postgresqlc.SetUserNotificationPreferencesParams{
		Secret:                  sqlc.XID(userSecret),
		NotificationPreferences: &prefs,
	})
}

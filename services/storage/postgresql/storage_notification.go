package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	"libdb.so/e2clicker/internal/sqlc/postgresqlc"
	"libdb.so/e2clicker/services/notification"
	"libdb.so/e2clicker/services/user"
)

func (s *Storage) notificationUserStorage() notification.UserNotificationStorage {
	return (*notificationUserStorage)(s)
}

type notificationUserStorage Storage

func (s *notificationUserStorage) UserPreferences(ctx context.Context, userSecret user.Secret) (notification.UserPreferences, error) {
	return s.q.UserNotificationPreferences(ctx, userSecret)
}

func (s *notificationUserStorage) SetUserPreferencesTx(ctx context.Context, userSecret user.Secret, prefs func(*notification.UserPreferences) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := postgresqlc.New(tx)

	p, err := q.UserNotificationPreferences(ctx, userSecret)
	if err != nil {
		return fmt.Errorf("get user preferences: %w", err)
	}

	if err := prefs(&p); err != nil {
		return err
	}

	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("cannot marshal UserPreferences as JSON: %w", err)
	}

	if err := q.SetUserNotificationPreferences(ctx, postgresqlc.SetUserNotificationPreferencesParams{
		Secret:  userSecret,
		Column2: b,
	}); err != nil {
		return fmt.Errorf("set user preferences: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

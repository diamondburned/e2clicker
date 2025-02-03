package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"e2clicker.app/internal/sqlc/postgresqlc"
	"e2clicker.app/services/user"
)

func (s *Storage) userStorage() user.UserStorage               { return s }
func (s *Storage) userSessionStorage() user.UserSessionStorage { return s }

func (s *Storage) CreateUser(ctx context.Context, userSecret user.Secret, name string) (user.User, error) {
	u, err := s.q.CreateUser(ctx, postgresqlc.CreateUserParams{
		Secret: userSecret,
		Name:   name,
	})
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		Name:   u.Name,
		Locale: u.Locale,
	}, nil
}

func (s *Storage) User(ctx context.Context, userSecret user.Secret) (user.User, error) {
	u, err := s.q.User(ctx, userSecret)
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		Name:   u.Name,
		Locale: u.Locale,
	}, nil
}

func (s *Storage) UpdateUserName(ctx context.Context, userSecret user.Secret, name string) error {
	return s.q.UpdateUserName(ctx, postgresqlc.UpdateUserNameParams{
		Secret: userSecret,
		Name:   name,
	})
}

func (s *Storage) UpdateUserLocale(ctx context.Context, userSecret user.Secret, locale user.Locale) error {
	return s.q.UpdateUserLocale(ctx, postgresqlc.UpdateUserLocaleParams{
		Secret: userSecret,
		Locale: locale,
	})
}

func (s *Storage) RegisterSession(ctx context.Context, token []byte, userSecret user.Secret, userAgent string) error {
	return s.q.RegisterSession(ctx, postgresqlc.RegisterSessionParams{
		UserSecret: userSecret,
		Token:      token,
		UserAgent:  pgtype.Text{String: userAgent, Valid: userAgent != ""},
	})
}

func (s *Storage) ValidateSession(ctx context.Context, token []byte) (user.Session, error) {
	r, err := s.q.ValidateSession(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.Session{}, user.ErrInvalidSession
		}
		return user.Session{}, err
	}
	return convertSession(r), nil
}

func (s *Storage) ListSessions(ctx context.Context, userSecret user.Secret) ([]user.Session, error) {
	l, err := s.q.ListSessions(ctx, userSecret)
	if err != nil {
		return nil, err
	}
	return convertList(l, convertSession), nil
}

func convertSession(r postgresqlc.UserSession) user.Session {
	return user.Session{
		ID:         r.ID,
		UserSecret: user.Secret(r.UserSecret),
		UserAgent:  r.UserAgent.String,
		CreatedAt:  r.CreatedAt.Time,
		LastUsed:   r.LastUsed.Time,
	}
}

func (s *Storage) DeleteSession(ctx context.Context, userSecret user.Secret, sessionID int64) error {
	err := s.q.DeleteSession(ctx, postgresqlc.DeleteSessionParams{
		UserSecret: userSecret,
		ID:         sessionID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return user.ErrInvalidSession
	}
	return err
}

func convertList[T1, T2 any](vs []T1, c func(T1) T2) []T2 {
	v2 := make([]T2, len(vs))
	for i, v := range vs {
		v2[i] = c(v)
	}
	return v2
}

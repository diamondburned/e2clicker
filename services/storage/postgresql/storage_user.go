package postgresql

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5/pgtype"
	"libdb.so/e2clicker/internal/asset"
	"libdb.so/e2clicker/internal/sqlc"
	"libdb.so/e2clicker/internal/sqlc/postgresqlc"
	"libdb.so/e2clicker/services/user"
)

func (s *Storage) userStorage() user.UserStorage               { return s }
func (s *Storage) userAvatarStorage() user.UserAvatarStorage   { return s }
func (s *Storage) userSessionStorage() user.UserSessionStorage { return s }

func (s *Storage) CreateUser(ctx context.Context, userSecret user.Secret, name string) (user.User, error) {
	u, err := s.q.CreateUser(ctx, postgresqlc.CreateUserParams{
		Secret: sqlc.XID(userSecret),
		Name:   name,
	})
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		Name:      u.Name,
		Locale:    u.Locale,
		HasAvatar: false,
	}, nil
}

func (s *Storage) User(ctx context.Context, userSecret user.Secret) (user.User, error) {
	u, err := s.q.User(ctx, sqlc.XID(userSecret))
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		Name:      u.Name,
		Locale:    u.Locale,
		HasAvatar: u.HasAvatar,
	}, nil
}

func (s *Storage) UpdateUserName(ctx context.Context, userSecret user.Secret, name string) error {
	return s.q.UpdateUserName(ctx, postgresqlc.UpdateUserNameParams{
		Secret: sqlc.XID(userSecret),
		Name:   name,
	})
}

func (s *Storage) UpdateUserLocale(ctx context.Context, userSecret user.Secret, locale user.Locale) error {
	return s.q.UpdateUserLocale(ctx, postgresqlc.UpdateUserLocaleParams{
		Secret: sqlc.XID(userSecret),
		Locale: locale,
	})
}

func (s *Storage) UserAvatar(ctx context.Context, userSecret user.Secret) (asset.ReadCloser, error) {
	a, err := s.q.UserAvatar(ctx, sqlc.XID(userSecret))
	if err != nil {
		return asset.ReadCloser{}, err
	}
	return asset.NewAssetReader(
		io.NopCloser(bytes.NewReader(a.AvatarImage)),
		a.MIMEType,
		int64(len(a.AvatarImage)),
	), nil
}

func (s *Storage) SetUserAvatar(ctx context.Context, id user.Secret, a asset.Reader) error {
	d, err := io.ReadAll(a.Reader())
	if err != nil {
		return fmt.Errorf("failed to read avatar: %w", err)
	}
	return s.q.SetUserAvatar(ctx, postgresqlc.SetUserAvatarParams{
		UserSecret:  sqlc.XID(id),
		MIMEType:    a.ContentType,
		AvatarImage: d,
	})
}

func (s *Storage) RegisterSession(ctx context.Context, token []byte, userSecret user.Secret, userAgent string) error {
	return s.q.RegisterSession(ctx, postgresqlc.RegisterSessionParams{
		UserSecret: sqlc.XID(userSecret),
		Token:      token,
		UserAgent:  pgtype.Text{String: userAgent, Valid: userAgent != ""},
	})
}

func (s *Storage) ValidateSession(ctx context.Context, token []byte) (user.Session, error) {
	r, err := s.q.ValidateSession(ctx, token)
	if err != nil {
		return user.Session{}, err
	}

	return user.Session{
		ID:         r.ID,
		UserSecret: user.Secret(r.UserSecret),
		UserAgent:  r.UserAgent.String,
		CreatedAt:  r.CreatedAt.Time,
		LastUsed:   r.LastUsed.Time,
		ExpiresAt:  r.ExpiresAt.Time,
	}, nil
}

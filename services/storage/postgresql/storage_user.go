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

var (
	_ user.UserStorage        = (*Storage)(nil)
	_ user.UserAvatarStorage  = (*Storage)(nil)
	_ user.UserSessionStorage = (*Storage)(nil)
)

func (s *Storage) CreateUser(ctx context.Context, id user.UserID, email string, passhash []byte, name string) (user.User, error) {
	u, err := s.q.CreateUser(ctx, postgresqlc.CreateUserParams{
		UserID:   sqlc.UserID{UserID: id},
		Email:    email,
		Passhash: passhash,
		Name:     name,
	})
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		ID:     u.UserID.UserID,
		Email:  u.Email,
		Name:   u.Name,
		Locale: u.Locale,
	}, nil
}

func (s *Storage) User(ctx context.Context, id user.UserID) (user.User, error) {
	u, err := s.q.User(ctx, sqlc.UserID{UserID: id})
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		ID:        u.UserID.UserID,
		Email:     u.Email,
		Name:      u.Name,
		Locale:    user.Locale(u.Locale),
		HasAvatar: u.HasAvatar,
	}, nil
}

func (s *Storage) UserPasswordFromEmail(ctx context.Context, email string) (user.UserPassword, error) {
	e, err := s.q.UserPasswordHashFromEmail(ctx, email)
	if err != nil {
		return user.UserPassword{}, err
	}
	return user.UserPassword{
		ID:       e.UserID.UserID,
		Passhash: e.Passhash,
	}, nil
}

func (s *Storage) UpdateUserEmailPassword(ctx context.Context, id user.UserID, email string, passhash []byte) error {
	return s.q.UpdateUserEmailPassword(ctx, postgresqlc.UpdateUserEmailPasswordParams{
		UserID:   sqlc.UserID{UserID: id},
		Email:    email,
		Passhash: passhash,
	})
}

func (s *Storage) UpdateUserName(ctx context.Context, id user.UserID, name string) error {
	return s.q.UpdateUserName(ctx, postgresqlc.UpdateUserNameParams{
		UserID: sqlc.UserID{UserID: id},
		Name:   name,
	})
}

func (s *Storage) UpdateUserLocale(ctx context.Context, id user.UserID, locale user.Locale) error {
	return s.q.UpdateUserLocale(ctx, postgresqlc.UpdateUserLocaleParams{
		UserID: sqlc.UserID{UserID: id},
		Locale: locale,
	})
}

func (s *Storage) UserAvatar(ctx context.Context, id user.UserID) (asset.ReadCloser, error) {
	a, err := s.q.UserAvatar(ctx, sqlc.UserID{UserID: id})
	if err != nil {
		return asset.ReadCloser{}, err
	}
	return asset.NewAssetReader(
		io.NopCloser(bytes.NewReader(a.AvatarImage)),
		a.MIMEType,
		int64(len(a.AvatarImage)),
	), nil
}

func (s *Storage) SetUserAvatar(ctx context.Context, id user.UserID, a asset.Reader) error {
	d, err := io.ReadAll(a.Reader())
	if err != nil {
		return fmt.Errorf("failed to read avatar: %w", err)
	}

	return s.q.SetUserAvatar(ctx, postgresqlc.SetUserAvatarParams{
		UserID:      sqlc.UserID{UserID: id},
		AvatarImage: d,
		MIMEType:    a.ContentType,
	})
}

func (s *Storage) RegisterSession(ctx context.Context, id user.UserID, token []byte, userAgent string) error {
	return s.q.RegisterSession(ctx, postgresqlc.RegisterSessionParams{
		UserID:    sqlc.UserID{UserID: id},
		Token:     token,
		UserAgent: pgtype.Text{String: userAgent, Valid: userAgent != ""},
	})
}

func (s *Storage) ValidateSession(ctx context.Context, token []byte) (user.User, error) {
	session, err := s.q.ValidateSession(ctx, token)
	if err != nil {
		return user.User{}, err
	}

	u, err := s.q.User(ctx, sqlc.UserID{UserID: session.UserID})
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:        u.UserID.UserID,
		Email:     u.Email,
		Name:      u.Name,
		Locale:    user.Locale(u.Locale),
		HasAvatar: u.HasAvatar,
	}, nil
}

package user

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"libdb.so/e2clicker/services/user/naming"
)

// UserService is a service for managing users.
type UserService struct {
	users        UserStorage
	userSessions UserSessionStorage
}

// UserServiceConfig is a dependency injection container for [UserService].
type UserServiceConfig struct {
	fx.In

	UserStorage
	UserSessionStorage
}

// NewUserService creates a new user service.
func NewUserService(c UserServiceConfig) (*UserService, error) {
	return &UserService{
		c.UserStorage,
		c.UserSessionStorage,
	}, nil
}

func (s UserService) CreateUser(ctx context.Context, name string) (UserWithSecret, error) {
	if name == "" {
		name = naming.RandomName()
	}
	secret := generateUserSecret()
	u, err := s.users.CreateUser(ctx, secret, name)
	if err != nil {
		return UserWithSecret{}, err
	}
	return UserWithSecret{u, secret}, nil
}

func (s UserService) User(ctx context.Context, secret Secret) (User, error) {
	return s.users.User(ctx, secret)
}

func (s UserService) UpdateUserName(ctx context.Context, secret Secret, name string) error {
	return s.users.UpdateUserName(ctx, secret, name)
}

func (s UserService) UpdateUserLocale(ctx context.Context, secret Secret, locale Locale) error {
	if err := locale.Validate(); err != nil {
		return fmt.Errorf("invalid locale: %w", err)
	}
	return s.users.UpdateUserLocale(ctx, secret, locale)
}

func (s UserService) CreateSession(ctx context.Context, userSecret Secret, userAgent string) (SessionToken, error) {
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	tokenBytes, err := token.asBytes()
	if err != nil {
		panic(err)
	}

	return token, s.userSessions.RegisterSession(ctx, tokenBytes, userSecret, userAgent)
}

func (s UserService) ValidateSession(ctx context.Context, token SessionToken) (Session, error) {
	tokenBytes, err := token.asBytes()
	if err != nil {
		return Session{}, ErrInvalidSession
	}

	session, err := s.userSessions.ValidateSession(ctx, tokenBytes)
	if err != nil {
		return Session{}, ErrInvalidSession
	}

	return session, nil
}

func (s UserService) ListSessions(ctx context.Context, userSecret Secret) ([]Session, error) {
	return s.userSessions.ListSessions(ctx, userSecret)
}

func (s UserService) DeleteSession(ctx context.Context, userSecret Secret, sessionID int64) error {
	return s.userSessions.DeleteSession(ctx, userSecret, sessionID)
}

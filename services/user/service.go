package user

import (
	"context"
	"fmt"
	"strings"

	scrypt "github.com/elithrar/simple-scrypt"
	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/asset"
)

// UserService is a service for managing users.
type UserService struct {
	users        UserStorage
	userAvatars  UserAvatarStorage
	userSessions UserSessionStorage
}

// UserServiceConfig is a dependency injection container for [UserService].
type UserServiceConfig struct {
	fx.In
	UserStorage
	UserAvatarStorage
	UserSessionStorage
}

// NewUserService creates a new user service.
func NewUserService(c UserServiceConfig) (*UserService, error) {
	return &UserService{
		c.UserStorage,
		c.UserAvatarStorage,
		c.UserSessionStorage,
	}, nil
}

func (s UserService) CreateUser(ctx context.Context, email, password, name string) (User, error) {
	if !isValidEmail(email) {
		return User{}, ErrInvalidEmail
	}

	id := GenerateUserID()

	passhash, err := hashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return s.users.CreateUser(ctx, id, email, passhash, name)
}

func hashPassword(password string) ([]byte, error) {
	return scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
}

func (s UserService) User(ctx context.Context, id UserID) (User, error) {
	return s.users.User(ctx, id)
}

func (s UserService) ValidateUserEmailPassword(ctx context.Context, email, password string) (UserID, error) {
	p, err := s.users.UserPasswordFromEmail(ctx, email)
	if err != nil {
		return NullUserID, ErrUnknownUser
	}

	if err := scrypt.CompareHashAndPassword(p.Passhash, []byte(password)); err != nil {
		return NullUserID, ErrUnknownUser
	}

	return p.ID, nil
}

func (s UserService) UpdateUserEmailPassword(ctx context.Context, id UserID, email, password string) error {
	if !isValidEmail(email) {
		return ErrInvalidEmail
	}

	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	passhash, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.users.UpdateUserEmailPassword(ctx, id, email, passhash)
}

func isValidEmail(email string) bool {
	return strings.Count(email, "@") == 1
}

func (s UserService) UpdateUserName(ctx context.Context, id UserID, name string) error {
	return s.users.UpdateUserName(ctx, id, name)
}

func (s UserService) UpdateUserLocale(ctx context.Context, id UserID, locale Locale) error {
	if err := locale.Validate(); err != nil {
		return fmt.Errorf("invalid locale: %w", err)
	}

	return s.users.UpdateUserLocale(ctx, id, locale)
}

func (s UserService) UserAvatar(ctx context.Context, id UserID) (asset.ReadCloser, error) {
	return s.userAvatars.UserAvatar(ctx, id)
}

func (s UserService) SetUserAvatar(ctx context.Context, id UserID, a asset.Reader) error {
	return s.userAvatars.SetUserAvatar(ctx, id, a)
}

func (s UserService) RegisterSession(ctx context.Context, id UserID, userAgent string) (SessionToken, error) {
	token, err := newSessionToken(id)
	if err != nil {
		return SessionToken{}, err
	}

	return token, s.userSessions.RegisterSession(ctx, id, token.randomBytes(), userAgent)
}

func (s UserService) ValidateSession(ctx context.Context, token string) (User, error) {
	parsed, err := ParseSessionToken(token)
	if err != nil {
		return User{}, ErrInvalidSession
	}

	user, err := s.userSessions.ValidateSession(ctx, parsed.randomBytes())
	if err != nil {
		return User{}, ErrInvalidSession
	}

	if parsed.UserID != user.ID {
		return User{}, ErrInvalidSession
	}

	return user, nil
}

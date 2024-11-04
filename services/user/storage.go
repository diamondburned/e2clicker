package user

import (
	"context"

	"libdb.so/e2clicker/internal/asset"
)

type UserStorage interface {
	// IsUserStorage is a marker method that ensures the interface is implemented.
	// It also prevents [UserService] from implementing this interface.
	// IsUserStorage()

	CreateUser(ctx context.Context, id UserID, email string, passhash []byte, name string) (User, error)
	User(ctx context.Context, id UserID) (User, error)
	UserPasswordFromEmail(ctx context.Context, email string) (UserPassword, error)
	UpdateUserEmailPassword(ctx context.Context, id UserID, email string, passhash []byte) error
	UpdateUserName(ctx context.Context, id UserID, name string) error
	UpdateUserLocale(ctx context.Context, id UserID, locale Locale) error
}

type UserAvatarStorage interface {
	// UserAvatar returns the user's avatar.
	// The returned asset is an open reader that must be closed by the caller.
	UserAvatar(ctx context.Context, id UserID) (asset.ReadCloser, error)
	// SetUserAvatar sets the user's avatar.
	// The caller must still close the given reader.
	SetUserAvatar(ctx context.Context, id UserID, a asset.Reader) error
}

type UserSessionStorage interface {
	// RegisterSession registers a session for a user. The token is generated by
	// [UserService]. The userAgent is optional.
	RegisterSession(ctx context.Context, id UserID, token []byte, userAgent string) error
	// ValidateSession validates a session for a user. The user that the session
	// belongs to is returned.
	ValidateSession(ctx context.Context, token []byte) (User, error)
}

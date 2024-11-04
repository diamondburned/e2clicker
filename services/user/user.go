// Package user provides services for managing users.
package user

import (
	"context"
	"encoding"
	"fmt"
	"time"

	"github.com/rs/xid"
	"libdb.so/e2clicker/internal/asset"
)

// User is a user in the system.
type User struct {
	Name      string
	Locale    Locale
	HasAvatar bool
}

// UserWithSecret is a user with their secret.
type UserWithSecret struct {
	User
	Secret Secret
}

type UserStorage interface {
	// CreateUser creates a user in the storage with the given name.
	CreateUser(ctx context.Context, secret Secret, name string) (User, error)
	// User gets the user identified by the given secret.
	User(ctx context.Context, secret Secret) (User, error)
	// UpdateUserName updates the user's name.
	UpdateUserName(ctx context.Context, secret Secret, name string) error
	// UpdateUserLocale updates the user's locale.
	UpdateUserLocale(ctx context.Context, secret Secret, locale Locale) error
}

type UserAvatarStorage interface {
	// UserAvatar returns the user's avatar.
	// The returned asset is an open reader that must be closed by the caller.
	UserAvatar(ctx context.Context, secret Secret) (asset.ReadCloser, error)
	// SetUserAvatar sets the user's avatar.
	// The caller must still close the given reader.
	SetUserAvatar(ctx context.Context, secret Secret, a asset.Reader) error
}

// Secret is a secret identifier for a user. This secret is generated once
// and never changes. It is used to both authenticate and identify a user, so it
// should be kept secret.
type Secret xid.ID

var (
	_ fmt.Stringer             = Secret{}
	_ encoding.TextMarshaler   = Secret{}
	_ encoding.TextUnmarshaler = (*Secret)(nil)
)

func generateUserSecret() Secret {
	return Secret(xid.New())
}

// String formats the user secret into a string.
func (id Secret) String() string {
	return xid.ID(id).String()
}

// CreatedAt returns the creation time of the user secret.
func (id Secret) CreatedAt() time.Time {
	return xid.ID(id).Time()
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (id Secret) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (id *Secret) UnmarshalText(text []byte) error {
	v, err := xid.FromString(string(text))
	if err != nil {
		return err
	}
	*id = Secret(v)
	return nil
}

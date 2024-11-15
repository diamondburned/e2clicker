// Package user provides services for managing users.
package user

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"strings"

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
type Secret string

func generateUserSecret() Secret {
	const n = 10

	var b [n]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}

	s := base32.
		StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(b[:])

	return Secret(s)
}

// String returns the secret as a pretty string.
func (s Secret) String() string {
	var b strings.Builder
	b.Grow(len(s) + (len(s)/4 + 1))
	for i, c := range s {
		if i != 0 && i%4 == 0 {
			b.WriteByte(' ')
		}
		b.WriteRune(c)
	}
	return b.String()
}

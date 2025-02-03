// Package user provides services for managing users.
package user

import (
	"context"
	"crypto/rand"
	"encoding"
	"encoding/base32"
	"strings"
)

// User is a user in the system.
type User struct {
	Name   string
	Locale Locale
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

// Secret is a secret identifier for a user. This secret is generated once
// and never changes. It is used to both authenticate and identify a user, so it
// should be kept secret.
type Secret string

var (
	_ encoding.TextMarshaler   = Secret("")
	_ encoding.TextUnmarshaler = (*Secret)(nil)
)

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

// PrettyString returns the secret as a pretty string.
func (s Secret) PrettyString() string {
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

func (s *Secret) UnmarshalText(text []byte) error {
	*s = Secret(text)
	*s = Secret(strings.ReplaceAll(string(*s), " ", ""))
	return nil
}

func (s Secret) MarshalText() ([]byte, error) {
	return []byte(string(s)), nil
}

// Package user provides services for managing users.
package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/text/language"
)

// UserID is a unique identifier for a user.
type UserID struct {
	x xid.ID
}

// NullUserID is a zero value for [UserID].
var NullUserID = UserID{}

// ParseUserID parses a user ID from a string. See [UserID.String].
func ParseUserID(s string) (UserID, error) {
	v, err := xid.FromString(s)
	if err != nil {
		return UserID{}, err
	}
	return UserID{v}, nil
}

// ParseRawUserID parses a user ID from raw bytes. See [UserID.Bytes].
func ParseRawUserID(b []byte) (UserID, error) {
	v, err := xid.FromBytes(b)
	if err != nil {
		return UserID{}, err
	}
	return UserID{v}, nil
}

// GenerateUserID generates a new user ID.
func GenerateUserID() UserID {
	return UserID{xid.New()}
}

// String formats the user ID into a string.
func (id UserID) String() string {
	return id.x.String()
}

// Bytes returns the user ID as raw bytes. This is useful for storing.
func (id UserID) Bytes() []byte {
	return id.x.Bytes()
}

// CreatedAt returns the creation time of the user ID.
func (id UserID) CreatedAt() time.Time {
	return id.x.Time()
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (id UserID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (id *UserID) UnmarshalText(text []byte) error {
	v, err := xid.FromString(string(text))
	if err != nil {
		return err
	}
	*id = UserID{v}
	return nil
}

// SessionToken is a token that represents a user session.
// The type represents an already validated token. The random part is not
// exposed to the user except via [String].
type SessionToken struct {
	UserID       UserID
	randomBase64 string
}

// ParseSessionToken parses a token string into a [SessionToken].
func ParseSessionToken(token string) (SessionToken, error) {
	rawUserID, randomBase64, ok := strings.Cut(token, ".")
	if !ok {
		return SessionToken{}, fmt.Errorf("invalid token format")
	}

	userID, err := xid.FromString(rawUserID)
	if err != nil {
		return SessionToken{}, fmt.Errorf("invalid user ID: %w", err)
	}

	return SessionToken{
		UserID:       UserID{userID},
		randomBase64: randomBase64,
	}, nil
}

func newSessionToken(userID UserID) (SessionToken, error) {
	var rawToken [24]byte

	_, err := io.ReadFull(rand.Reader, rawToken[:])
	if err != nil {
		return SessionToken{}, fmt.Errorf("failed to generate session token: %w", err)
	}

	return SessionToken{
		UserID:       userID,
		randomBase64: base64.URLEncoding.EncodeToString(rawToken[:]),
	}, nil
}

// String returns the token as a string.
func (t SessionToken) String() string { return t.UserID.String() + "." + t.randomBase64 }

func (t SessionToken) randomBytes() []byte {
	b, err := base64.URLEncoding.DecodeString(t.randomBase64)
	if err != nil {
		panic("invalid base64 in session token (bug; only use this method on new tokens)")
	}
	return b
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (t SessionToken) MarshalText() ([]byte, error) { return []byte(t.String()), nil }

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (t *SessionToken) UnmarshalText(text []byte) error {
	token, err := ParseSessionToken(string(text))
	if err != nil {
		return err
	}
	*t = token
	return nil
}

// User is a user in the system.
type User struct {
	ID        UserID `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Locale    Locale `json:"locale"`
	HasAvatar bool   `json:"hasAvatar,omitempty"`
}

// Locale is a user's preferred languages. It is used for localization.
// The format of the string is specified by RFC 2616 but is validated by
// [language.ParseAcceptLanguage], which is more lax.
type Locale string

// ParseLocale parses a locale string into a [Locale] type.
func ParseLocale(locale string) (Locale, error) {
	l := Locale(locale)
	return l, l.Validate()
}

// Tags returns the Locale as a list of language tags.
// If l is empty or invalid, then this function returns one [language.Und]. The
// returned list is never empty.
func (l Locale) Tags() []language.Tag {
	tags, _, _ := language.ParseAcceptLanguage(string(l))
	if len(tags) == 0 {
		return []language.Tag{language.Und}
	}
	return tags
}

// Validate checks if the Locale is valid.
func (l Locale) Validate() error {
	_, _, err := language.ParseAcceptLanguage(string(l))
	return err
}

// String implements the [fmt.Stringer] interface.
func (l Locale) String() string {
	return string(l)
}

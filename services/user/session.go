package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

type UserSessionStorage interface {
	// RegisterSession registers a session for a user. The token is generated by
	// [UserService]. The userAgent is optional.
	RegisterSession(ctx context.Context, token []byte, userSecret Secret, userAgent string) error
	// ValidateSession validates a session for a user. The user that the session
	// belongs to is returned.
	ValidateSession(ctx context.Context, token []byte) (Session, error)
	// ListSessions lists all sessions for a user.
	ListSessions(ctx context.Context, userSecret Secret) ([]Session, error)
	// DeleteSession deletes a session for a user.
	DeleteSession(ctx context.Context, userSecret Secret, sessionID int64) error
}

// Session is a user session.
type Session struct {
	// ID uniquely identifies the session.
	ID int64
	// UserSecret is the secret of the user that the session belongs to.
	UserSecret Secret
	// UserAgent is the user agent that the session was created with.
	UserAgent string
	// CreatedAt is the time that the session was created.
	CreatedAt time.Time
	// LastUsed is the time that the session was last used.
	LastUsed time.Time
	// ExpiresAt is the time that the session will expire.
	// If zero, the session does not expire.
	ExpiresAt time.Time
}

// SessionToken is a token that represents a user session.
// The type represents an already validated token. The random part is not
// exposed to the user except via [String].
type SessionToken string

// sessionTokenFromBytes creates a new session token from raw bytes.
func sessionTokenFromBytes(b []byte) SessionToken {
	return SessionToken(base64.URLEncoding.EncodeToString(b))
}

// generateSessionToken generates a new session token.
func generateSessionToken() (SessionToken, error) {
	var rawToken [24]byte

	_, err := io.ReadFull(rand.Reader, rawToken[:])
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	return sessionTokenFromBytes(rawToken[:]), nil
}

// asBytes returns the session token as raw bytes.
func (t SessionToken) asBytes() ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(string(t))
	if err != nil {
		return nil, ErrInvalidSession
	}
	return b, nil
}

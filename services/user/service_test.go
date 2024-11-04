package user

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()

	s := newMockUserService(t)

	_, err := s.CreateUser(ctx, "Diamond")
	assert.NoError(t, err)

	call := s.users.CreateUserCalls()[0]
	assert.NotZero(t, call.Secret)
	assert.Equal(t, call.Name, "Diamond")
}

func TestUserService_UpdateUserLocale(t *testing.T) {
	tests := []struct {
		locale Locale
		valid  bool
	}{
		{"en", true},
		{"fr", true},
		{"en-US", true},
		{"fr-CA, fr;q=0.9, en;q=0.8", true},
		{"whatever lol", false},
	}

	ctx := context.Background()
	secret := generateUserSecret()

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s := newMockUserService(t)

			err := s.UpdateUserLocale(ctx, secret, test.locale)
			if !test.valid {
				assert.Error(t, err)
				assert.Equal(t, len(s.users.UpdateUserLocaleCalls()), 0)
				return
			}

			assert.NoError(t, err)

			call := s.users.UpdateUserLocaleCalls()[0]
			assert.Equal(t, call.Secret, secret)
			assert.Equal(t, call.Locale, test.locale)
		})
	}
}

func TestUserService_CreateSession(t *testing.T) {
	ctx := context.Background()
	secret := generateUserSecret()

	s := newMockUserService(t)

	token, err := s.CreateSession(ctx, secret, "user agent")
	assert.NoError(t, err)
	assert.NotZero(t, token)

	tokenBytes, err := token.asBytes()
	assert.NoError(t, err)

	s.sessions.ValidateSessionFunc = func(ctx context.Context, wantToken []byte) (Session, error) {
		if !bytes.Equal(wantToken, tokenBytes) {
			return Session{}, fmt.Errorf("token not found")
		}
		return Session{
			ID:         1,
			UserSecret: secret,
			UserAgent:  "user agent",
		}, nil
	}

	t.Run("correct", func(t *testing.T) {
		s, err := s.ValidateSession(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, s.ID, int64(1))
		assert.Equal(t, s.UserSecret, secret)
		assert.Equal(t, s.UserAgent, "user agent")
	})

	t.Run("incorrect", func(t *testing.T) {
		s, err := s.ValidateSession(ctx, "wrong token")
		assert.Equal(t, err, ErrInvalidSession)
		assert.Zero(t, s)
	})

	t.Run("tampered", func(t *testing.T) {
		token2 := token
		token2 = token2[:len(token2)-1] + "0"

		s, err := s.ValidateSession(ctx, token2)
		assert.Equal(t, err, ErrInvalidSession)
		assert.Zero(t, s)
	})
}

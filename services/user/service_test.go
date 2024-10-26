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

	_, err := s.CreateUser(ctx, "diamond@example.com", "password", "Diamond")
	assert.NoError(t, err)

	call := s.users.CreateUserCalls()[0]
	assert.NotZero(t, call.ID)
	assert.Equal(t, call.Name, "Diamond")
	assert.Equal(t, call.Email, "diamond@example.com")
}

func TestUserService_ValidateUserEmailPassword(t *testing.T) {
	const rightEmail = "otheremail@example.com"
	const rightPassword = "password"

	ctx := context.Background()
	rightID := GenerateUserID()

	s := newMockUserService(t)
	s.users.UserPasswordFromEmailFunc = func(ctx context.Context, email string) (UserPassword, error) {
		if email == rightEmail {
			b, _ := hashPassword(rightPassword)
			return UserPassword{ID: rightID, Passhash: b}, nil
		}
		return UserPassword{}, fmt.Errorf("user not found")
	}

	var id UserID
	var err error

	id, err = s.ValidateUserEmailPassword(ctx, rightEmail, rightPassword)
	assert.NoError(t, err)
	assert.Equal(t, id, rightID)

	id, err = s.ValidateUserEmailPassword(ctx, "wrongemail", rightPassword)
	assert.Equal(t, err, ErrUnknownUser)
	assert.Equal(t, id, NullUserID)

	id, err = s.ValidateUserEmailPassword(ctx, rightEmail, "wrongpassword")
	assert.Equal(t, err, ErrUnknownUser)
	assert.Equal(t, id, NullUserID)
}

func TestUserService_UpdateUserEmailPassword(t *testing.T) {
	ctx := context.Background()
	id := GenerateUserID()

	s := newMockUserService(t)

	err := s.UpdateUserEmailPassword(ctx, id, "otheremail@example.com", "password2")
	assert.NoError(t, err)

	call := s.users.UpdateUserEmailPasswordCalls()[0]
	assert.Equal(t, call.ID, id)
	assert.Equal(t, call.Email, "otheremail@example.com")
	assert.NotEqual(t, string(call.Passhash), "password2")

	err = s.UpdateUserEmailPassword(ctx, NullUserID, "wrongemail", "password2")
	assert.Equal(t, err, ErrInvalidEmail)

	err = s.UpdateUserEmailPassword(ctx, NullUserID, "otheremail@example.com", "")
	assert.Equal(t, err, ErrPasswordTooShort)
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
	id := GenerateUserID()

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s := newMockUserService(t)

			err := s.UpdateUserLocale(ctx, id, test.locale)
			if !test.valid {
				assert.Error(t, err)
				assert.Equal(t, len(s.users.UpdateUserLocaleCalls()), 0)
				return
			}

			assert.NoError(t, err)

			call := s.users.UpdateUserLocaleCalls()[0]
			assert.Equal(t, call.ID, id)
			assert.Equal(t, call.Locale, test.locale)
		})
	}
}

func TestUserService_RegisterSession(t *testing.T) {
	ctx := context.Background()
	id := GenerateUserID()

	s := newMockUserService(t)

	token, err := s.RegisterSession(ctx, id, "user agent")
	assert.NoError(t, err)
	assert.Equal(t, token.UserID, id)

	s.sessions.ValidateSessionFunc = func(ctx context.Context, wantToken []byte) (User, error) {
		if !bytes.Equal(wantToken, token.randomBytes()) {
			return User{}, fmt.Errorf("token not found")
		}
		return User{ID: id}, nil
	}

	t.Run("correct", func(t *testing.T) {
		u, err := s.ValidateSession(ctx, token.String())
		assert.NoError(t, err)
		assert.Equal(t, u.ID, id)
	})

	t.Run("incorrect", func(t *testing.T) {
		_, err := s.ValidateSession(ctx, "wrong token")
		assert.Equal(t, err, ErrInvalidSession)
	})

	t.Run("tampered", func(t *testing.T) {
		token2 := token
		token2.UserID = GenerateUserID()

		_, err := s.ValidateSession(ctx, token2.String())
		assert.Equal(t, err, ErrInvalidSession)
	})
}

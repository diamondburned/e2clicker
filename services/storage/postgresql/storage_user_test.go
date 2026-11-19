package postgresql

import (
	"context"
	"errors"
	"testing"

	e2clickermodule "e2clicker.app/nix/modules/e2clicker"
	"e2clicker.app/services/user"
	"github.com/alecthomas/assert/v2"
	"github.com/neilotoole/slogt"
	"go.uber.org/fx/fxtest"
)

func TestUserAndSessionUnknownUser(t *testing.T) {
	ctx := context.Background()
	uri := startTestPostgres(t)

	lc := fxtest.NewLifecycle(t)
	s, err := NewStorage(lc, e2clickermodule.PostgreSQL{DatabaseURI: uri}, slogt.New(t))
	assert.NoError(t, err)
	lc.RequireStart()
	defer lc.RequireStop()

	users := s.userStorage()
	sessions := s.userSessionStorage()

	const missingSecret user.Secret = "NONEXISTENTUSER1"

	_, err = users.User(ctx, missingSecret)
	assert.True(t, errors.Is(err, user.ErrUnknownUser))

	err = sessions.RegisterSession(ctx, []byte("token-1"), missingSecret, "test-agent")
	assert.True(t, errors.Is(err, user.ErrUnknownUser))

	const validSecret user.Secret = "EXISTENTUSER0001"
	_, err = users.CreateUser(ctx, validSecret, "Alice")
	assert.NoError(t, err)

	err = sessions.RegisterSession(ctx, []byte("token-2"), validSecret, "test-agent")
	assert.NoError(t, err)
}

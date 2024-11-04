package api

import (
	"context"
	"errors"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"libdb.so/ctxt"
	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/e2clicker/services/user"
)

func init() {
	publicerrors.MarkValuesPublic(ErrNotBearerAuth)
}

var ErrNotBearerAuth = errors.New("not a bearer authentication token")

// Authenticator is an authenticator that authenticates requests.
type Authenticator struct {
	users *user.UserService
}

// NewAuthenticator creates a new authenticator.
func NewAuthenticator(users *user.UserService) *Authenticator {
	return &Authenticator{
		users,
	}
}

func (a *Authenticator) AuthenticationFunc() openapi3filter.AuthenticationFunc {
	return a.authenticate
}

func (a *Authenticator) authenticate(ctx context.Context, auth *openapi3filter.AuthenticationInput) error {
	if auth.SecuritySchemeName == "BearerAuth" {
		return a.authenticateBearerAuth(ctx, auth)
	}
	return auth.NewError(ErrNotBearerAuth)
}

func (a *Authenticator) authenticateBearerAuth(ctx context.Context, auth *openapi3filter.AuthenticationInput) error {
	r := auth.RequestValidationInput.Request

	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		return auth.NewError(ErrNotBearerAuth)
	}

	s, err := a.users.ValidateSession(ctx, user.SessionToken(token))
	if err != nil {
		return auth.NewError(err)
	}

	// Awful hack to pass the session to the handler.
	// Blame the OpenAPI generator for this.
	*auth.RequestValidationInput.Request = *r.WithContext(ctxt.With(ctx, s))

	return nil
}

func sessionFromCtx(ctx context.Context) user.Session {
	s, ok := ctxt.From[user.Session](ctx)
	if ok {
		return s
	}
	panic("BUG: session not found in context")
}

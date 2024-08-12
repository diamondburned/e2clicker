package user

import (
	"errors"

	"libdb.so/e2clicker/internal/publicerrors"
)

func init() {
	publicerrors.MarkValuesPublic(
		ErrInvalidEmail,
		ErrUnknownUser,
		ErrPasswordTooShort,
		ErrInvalidSession,
	)
}

// ErrInvalidEmail is returned when the email is invalid.
var ErrInvalidEmail = errors.New("invalid email")

// ErrUnknownUser is returned when the user is unknown.
var ErrUnknownUser = errors.New("unknown user")

// ErrPasswordTooShort is returned when the password is too short.
var ErrPasswordTooShort = errors.New("password too short")

// ErrInvalidSession is returned when the session is invalid, either because it
// is unknown or expired.
var ErrInvalidSession = errors.New("invalid session")

package postgresqlc

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeTableNotFound   = "42P01"
	CodeUniqueViolation = "23505"
)

// IsConstraintFailed returns true if err is returned because of a unique
// constraint violation.
func IsConstraintFailed(err error) bool {
	return IsErrorCode(err, CodeUniqueViolation)
}

// IsErrorCode returns true if err is a PostgreSQL error with the given code.
func IsErrorCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	// https://www.postgresql.org/docs/current/errcodes-appendix.html
	return pgErr.Code == code
}

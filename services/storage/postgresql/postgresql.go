// Package postgresql provides an implementation of the storage service using
// PostgreSQL. It relies on the generated sqlc queries.
package postgresql

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"libdb.so/e2clicker/services/storage/sqlc/postgresqlc"
)

const (
	codeTableNotFound   = "42P01"
	codeUniqueViolation = "23505"
)

// Connect connects to a pgSQL database.
func Connect(ctx context.Context, url string) (*Storage, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse url")
	}

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create pool")
	}

	if err := migrate(ctx, conn); err != nil {
		return nil, errors.Wrap(err, "cannot migrate")
	}

	return newStorage(conn), nil
}

func migrate(ctx context.Context, conn *pgxpool.Pool) error {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable, // force strictest isolation level
	})
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}
	defer tx.Rollback(ctx)

	v, err := postgresqlc.New(tx).Version(ctx)
	if err != nil {
		if !isErrorCode(err, codeTableNotFound) {
			return errors.Wrap(err, "cannot get schema version")
		}
	}

	versions := postgresqlc.Schema.Versions()
	if int(v) > len(versions) {
		return fmt.Errorf(
			"database schema version %d is higher than the latest supported version %d (app outdated?)",
			v, len(versions))
	}

	if int(v) == len(versions) {
		return nil
	}

	for i := int(v); i < len(versions); i++ {
		_, err := tx.Exec(ctx, versions[i])
		if err != nil {
			return errors.Wrapf(err, "cannot apply migration %d (from 0th)", i)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "cannot commit new migrations")
	}

	return nil
}

// isConstraintFailed returns true if err is returned because of a unique
// constraint violation.
func isConstraintFailed(err error) bool {
	return isErrorCode(err, codeUniqueViolation)
}

func isErrorCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	// https://www.postgresql.org/docs/current/errcodes-appendix.html
	return pgErr.Code == code
}

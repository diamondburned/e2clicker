// Package postgresql provides an implementation of the storage service using
// PostgreSQL. It relies on the generated sqlc queries.
package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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
		return nil, fmt.Errorf("cannot parse url: %w", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create pool: %w", err)
	}

	if err := migrate(ctx, conn); err != nil {
		return nil, fmt.Errorf("cannot migrate: %w", err)
	}

	return newStorage(conn), nil
}

func migrate(ctx context.Context, conn *pgxpool.Pool) error {
	var firstRun bool

	currentVersion, err := postgresqlc.New(conn).Version(ctx)
	if err != nil {
		if !isErrorCode(err, codeTableNotFound) {
			return fmt.Errorf("cannot get schema version: %w", err)
		}
		firstRun = true
	}

	schemaVersions := postgresqlc.Schema.Versions()
	if int(currentVersion) >= len(schemaVersions) {
		return nil
	}

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable, // force strictest isolation level
	})
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if !firstRun {
		currentVersion, err = postgresqlc.New(tx).Version(ctx)
		if err != nil {
			return fmt.Errorf("cannot get schema version: %w", err)
		}
	}

	for i := int(currentVersion); i < len(schemaVersions); i++ {
		_, err := tx.Exec(ctx, schemaVersions[i])
		if err != nil {
			slog.Error(
				"cannot apply migration",
				"module", "storage.postgresql",
				"version", i,
				"version_initial", currentVersion,
				"version_wanted", len(schemaVersions),
				"err", err)
			return fmt.Errorf("cannot apply migration %d", i)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit new migrations: %w", err)
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

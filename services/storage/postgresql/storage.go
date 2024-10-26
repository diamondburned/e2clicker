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
	"github.com/samber/do/v2"
	"libdb.so/e2clicker/internal/sqlc/postgresqlc"
)

// Config is the configuration for the storage service.
type Config struct {
	DatabaseURI string `json:"databaseURI"`
}

func (c *Config) Validate() error {
	if c.DatabaseURI == "" {
		return errors.New("database_uri is required")
	}
	return nil
}

const (
	codeTableNotFound   = "42P01"
	codeUniqueViolation = "23505"
)

// Storage is the PostgreSQL-backed storage.
type Storage struct {
	q       *postgresqlc.Queries
	conn    *pgxpool.Pool
	conncfg *pgxpool.Config
	logger  *slog.Logger
}

var (
	_ do.ShutdownerWithContextAndError = (*Storage)(nil)
	_ do.HealthcheckerWithContext      = (*Storage)(nil)
)

func newStorage(i do.Injector) (*Storage, error) {
	ctx := do.MustInvoke[context.Context](i)
	logger := do.MustInvoke[*slog.Logger](i)
	config := do.MustInvoke[*Config](i)

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	conncfg, err := pgxpool.ParseConfig(config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("cannot parse url: %w", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, conncfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create pool: %w", err)
	}

	if err := migrate(ctx, conn); err != nil {
		return nil, fmt.Errorf("cannot migrate: %w", err)
	}

	return &Storage{
		postgresqlc.New(conn),
		conn,
		conncfg,
		logger,
	}, nil
}

func (s *Storage) HealthCheck(ctx context.Context) error {
	return s.conn.Ping(ctx)
}

func (s *Storage) Shutdown(ctx context.Context) error {
	s.conn.Close()
	return nil
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

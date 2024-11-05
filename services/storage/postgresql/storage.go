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
	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/sqlc/postgresqlc"

	e2clickermodule "libdb.so/e2clicker/nix/modules/e2clicker"
)

const (
	codeTableNotFound   = "42P01"
	codeUniqueViolation = "23505"
)

// Storage is the PostgreSQL-backed storage.
type Storage struct {
	q       *postgresqlc.Queries // nil until start
	pool    *pgxpool.Pool        // nil until start
	conncfg *pgxpool.Config
	logger  *slog.Logger
}

// NewStorage creates a new PostgreSQL-backed storage.
func NewStorage(lc fx.Lifecycle, config e2clickermodule.PostgreSQL, logger *slog.Logger) (*Storage, error) {
	if config.DatabaseURI == "" {
		return nil, errors.New("database_uri is required")
	}

	conncfg, err := pgxpool.ParseConfig(config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("cannot parse url: %w", err)
	}

	logger = logger.With(slog.Group(
		"postgresql",
		"host", conncfg.ConnConfig.Host,
		"database", conncfg.ConnConfig.Database,
	))

	s := &Storage{
		q:       nil,
		pool:    nil,
		conncfg: conncfg,
		logger:  logger,
	}

	poolStatLogs := func() slog.Attr {
		if s.pool == nil {
			return slog.Attr{}
		}
		poolStat := s.pool.Stat()
		return slog.Group("pool_stat",
			"total_conns", poolStat.TotalConns(),
			"idle_conns", poolStat.IdleConns())
	}

	conncfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		logger.DebugContext(ctx, "adding new PostgreSQL connection to pool", poolStatLogs())
		return nil
	}

	conncfg.BeforeClose = func(conn *pgx.Conn) {
		logger.Debug("closing PostgreSQL connection in pool", poolStatLogs())
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pool, err := pgxpool.NewWithConfig(ctx, conncfg)
			if err != nil {
				return fmt.Errorf("cannot create pool: %w", err)
			}

			if err := migrate(ctx, pool); err != nil {
				return fmt.Errorf("cannot migrate: %w", err)
			}

			s.q = postgresqlc.New(pool)
			s.pool = pool

			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.pool.Close()
			return nil
		},
	})

	return s, nil
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

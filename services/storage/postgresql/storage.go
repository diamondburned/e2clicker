// Package postgresql provides an implementation of the storage service using
// PostgreSQL. It relies on the generated sqlc queries.
package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"libdb.so/e2clicker/internal/sqlc/postgresqlc"

	e2clickermodule "libdb.so/e2clicker/nix/modules/e2clicker"
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

			if err := postgresqlc.Migrate(ctx, pool); err != nil {
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

func maybePtr[T any](v T, valid bool) *T {
	if !valid {
		return nil
	}
	return &v
}

func deref[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// Package postgresql provides an implementation of the storage service using
// PostgreSQL. It relies on the generated sqlc queries.
package postgresql

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"

	"e2clicker.app/internal/slogutil"
	"e2clicker.app/internal/sqlc/postgresqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	e2clickermodule "e2clicker.app/nix/modules/e2clicker"
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

type copyFromIterator[T any] struct {
	value T
	error error

	next func() (T, bool)
	stop func()
	rows func(T) ([]any, error)
}

var (
	_ pgx.CopyFromSource = (*copyFromIterator[struct{}])(nil)
	_ io.Closer          = (*copyFromIterator[struct{}])(nil)
)

func newCopyFromIterator[T any](seq iter.Seq[T], rows func(T) ([]any, error)) *copyFromIterator[T] {
	// Pull's recover() isn't particularly helpful, so we wrap it.
	next, stop := iter.Pull(func(yield func(T) bool) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error(
					"panic in iterator used by storage/postgresql.copyFromIterator",
					"panic", r,
					slogutil.StackTrace(1))
				panic(r)
			}
		}()

		for v := range seq {
			if !yield(v) {
				return
			}
		}
	})

	return &copyFromIterator[T]{
		next: next,
		stop: stop,
		rows: rows,
	}
}

func (i *copyFromIterator[T]) Next() (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false

			slog.Error(
				"panic in storage/postgresql.copyFromIterator",
				"panic", r,
				slogutil.StackTrace(1))
		}
	}()
	i.value, ok = i.next()
	return
}

func (i *copyFromIterator[T]) Values() ([]any, error) {
	return i.rows(i.value)
}

func (i *copyFromIterator[T]) Err() error {
	return nil
}

func (i *copyFromIterator[T]) Close() error {
	i.stop()
	return nil
}

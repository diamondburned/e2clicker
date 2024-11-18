// Package postgresqlc contains raw SQL queries for the storage services and
// wrapper functions to use them.
package postgresqlc

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"libdb.so/lazymigrate"
)

//go:embed schema/*.sql
var schemaFS embed.FS

// Migrate runs all migrations on the database.
func Migrate(ctx context.Context, conn *pgxpool.Pool) error {
	return migrateAll(ctx, conn)
}

func migrateAll(ctx context.Context, conn *pgxpool.Pool) error {
	// force strictest isolation level
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	entries := must2(fs.ReadDir(schemaFS, "schema"))
	// Sort alphabetically. Our entries are prefixed with NN-, so a string
	// comparison is enough.
	slices.SortFunc(entries, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	for _, entry := range entries {
		data := must2(fs.ReadFile(schemaFS, "schema/"+entry.Name()))
		isPrimary := strings.HasSuffix(entry.Name(), "-primary.sql")

		slog.DebugContext(ctx,
			"postgresqlc: migrating schema",
			"schema", entry.Name(),
			"is_primary", isPrimary)

		var err error
		if isPrimary {
			schema := lazymigrate.NewSchemaWithMagic(string(data), "-- NEW VERSION")
			err = migratePrimary(ctx, tx, schema)
		} else {
			_, err = tx.Exec(ctx, string(data))
		}

		if err != nil {
			slog.ErrorContext(ctx,
				"cannot apply migration",
				"schema", entry.Name(),
				"is_primary", isPrimary,
				tint.Err(err))
			return fmt.Errorf("cannot apply migration %s", entry.Name())
		}

		slog.DebugContext(ctx,
			"postgresqlc: migrated schema",
			"schema", entry.Name(),
			"is_primary", isPrimary)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit new migrations: %w", err)
	}

	return nil
}

func migratePrimary(ctx context.Context, tx0 pgx.Tx, schema *lazymigrate.Schema) error {
	tx, err := tx0.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	vlatest := len(schema.Versions()) + 1

	for v := 1; v < vlatest; v++ {
		_, err := tx.Exec(ctx, schema.Versions()[v-1])
		if err != nil {
			slog.Error(
				"cannot apply migration",
				"module", "storage.postgresql",
				"version", v,
				"version_initial", v,
				"version_wanted", vlatest,
				"err", err)
			return fmt.Errorf("cannot apply version %d", v)
		}

		// Get the current version and set v to it.
		// Afterwards, v will always be incremented to the next number before
		// the next schema migration is applied, if any.
		c, err := New(tx).Version(ctx)
		if err != nil {
			return fmt.Errorf("cannot get version: %w", err)
		}
		v = int(c)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit new versions: %w", err)
	}

	return nil
}

func must2[T any](v T, err error) T {
	must1(err)
	return v
}

func must1(err error) {
	if err != nil {
		panic(err)
	}
}

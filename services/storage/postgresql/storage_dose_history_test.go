package postgresql

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
	"time"
	"uuid"

	e2clickermodule "e2clicker.app/nix/modules/e2clicker"
	"e2clicker.app/services/dosage"
	"github.com/alecthomas/assert/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neilotoole/slogt"
	"go.uber.org/fx/fxtest"
	"libdb.so/lazymigrate"
)

func startTestPostgres(t *testing.T) string {
	t.Helper()

	for _, bin := range []string{"initdb", "pg_ctl", "createdb"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not found in PATH: %v", bin, err)
		}
	}

	work, err := os.MkdirTemp("", "e2pg*")
	assert.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(work) })

	pgdata := filepath.Join(work, "data")
	pgsock := filepath.Join(work, "sock")
	assert.NoError(t, os.MkdirAll(pgdata, 0o700))
	assert.NoError(t, os.MkdirAll(pgsock, 0o700))

	run := func(name string, args ...string) {
		cmd := exec.Command(name, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
		}
	}

	run("initdb", "-D", pgdata, "--username=postgres", "--auth=trust", "--no-sync")
	run("pg_ctl", "-D", pgdata, "-l", filepath.Join(work, "pg.log"),
		"-o", fmt.Sprintf("-k %s -h ''", pgsock), "start")
	t.Cleanup(func() {
		_ = exec.Command("pg_ctl", "-D", pgdata, "-m", "immediate", "stop").Run()
	})

	run("createdb", "-h", pgsock, "-U", "postgres", "e2clicker")

	return "postgres://postgres@/e2clicker?host=" + pgsock
}

func TestDoseHistoryMigrationAndStorage(t *testing.T) {
	ctx := context.Background()
	uri := startTestPostgres(t)

	// Phase 1: initialize schema at v5, seed data (including a pre-1970 dose),
	// then run postgresqlc.Migrate to verify v5 -> v6 upgrade.
	pool, err := pgxpool.New(ctx, uri)
	assert.NoError(t, err)
	defer pool.Close()

	deliveryMethodsSQL, err := os.ReadFile(filepath.Join("..", "..", "..", "internal", "sqlc", "postgresqlc", "schema", "01-delivery_methods.sql"))
	assert.NoError(t, err)
	_, err = pool.Exec(ctx, string(deliveryMethodsSQL))
	assert.NoError(t, err)

	primarySQL, err := os.ReadFile(filepath.Join("..", "..", "..", "internal", "sqlc", "postgresqlc", "schema", "99-primary.sql"))
	assert.NoError(t, err)

	schema := lazymigrate.NewSchemaWithMagic(string(primarySQL), "-- NEW VERSION")
	versions := schema.Versions()
	assert.True(t, len(versions) >= 6, "expected at least 6 schema versions")

	for i := range 5 {
		_, err := pool.Exec(ctx, versions[i])
		assert.NoError(t, err)
	}

	var version int16
	assert.NoError(t, pool.QueryRow(ctx, "SELECT v FROM meta").Scan(&version))
	assert.Equal(t, int16(5), version)

	_, err = pool.Exec(ctx, `INSERT INTO users (secret, name) VALUES ('u1', 'User 1'), ('u2', 'User 2')`)
	assert.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO dosage_history (user_secret, delivery_method, dose, taken_at, comment) VALUES
			('u1', 'EV im', 5, '2025-01-10T12:00:00Z', 'migrated 1'),
			('u1', 'EV im', 4, '1965-01-10T12:00:00Z', 'pre-epoch')
	`)
	assert.NoError(t, err)

	// Close initial pool before starting Storage (which runs Migrate).
	pool.Close()

	lc := fxtest.NewLifecycle(t)
	s, err := NewStorage(lc, e2clickermodule.PostgreSQL{DatabaseURI: uri}, slogt.New(t))
	assert.NoError(t, err)
	lc.RequireStart()
	defer lc.RequireStop()

	assert.NoError(t, s.pool.QueryRow(ctx, "SELECT v FROM meta").Scan(&version))
	assert.Equal(t, int16(6), version)

	// Verify migrated rows have valid RFC 9562 UUIDv7 IDs.
	hist := s.doseHistoryStorage()
	u1History, err := hist.DoseHistory(ctx, "u1", time.Time{}, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(u1History))
	for _, d := range u1History {
		assert.NotEqual(t, uuid.Nil(), d.ID)
		assert.Equal(t, byte(7), d.ID[6]>>4)
		assert.Equal(t, byte(2), d.ID[8]>>6)
	}

	// Verify old unique constraint / indexes were dropped and new index exists.
	var oldIdxCount int
	assert.NoError(t, s.pool.QueryRow(ctx, `
		SELECT count(*) FROM pg_indexes
		WHERE indexname IN ('dosage_history_user_secret_taken_at_key', 'dosage_history_taken_at', 'dosage_history_user_secret')
	`).Scan(&oldIdxCount))
	assert.Equal(t, 0, oldIdxCount)

	var newIdxCount int
	assert.NoError(t, s.pool.QueryRow(ctx, `
		SELECT count(*) FROM pg_indexes
		WHERE indexname IN ('dosage_history_user_taken_at_id', 'dosage_history_pkey')
	`).Scan(&newIdxCount))
	assert.Equal(t, 2, newIdxCount)

	// Phase 2: same-second RecordDose, EditDose by UUID, ForgetDoses by UUID.
	sameTime := time.Date(2025, 2, 1, 10, 0, 0, 0, time.UTC)
	d1, err := hist.RecordDose(ctx, "u1", dosage.Dose{
		DeliveryMethod: "EV im",
		Dose:           5,
		TakenAt:        sameTime,
		Comment:        "first",
	})
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil(), d1.ID)

	d2, err := hist.RecordDose(ctx, "u1", dosage.Dose{
		DeliveryMethod: "EV im",
		Dose:           5,
		TakenAt:        sameTime,
		Comment:        "second",
	})
	assert.NoError(t, err)
	assert.NotEqual(t, d1.ID, d2.ID)

	// Edit only d1; verify d2 is untouched and comment persists.
	assert.NoError(t, hist.EditDose(ctx, "u1", d1.ID, dosage.Dose{
		DeliveryMethod: "EV im",
		Dose:           6,
		TakenAt:        sameTime,
		Comment:        "edited first",
	}))

	u1AfterEdit, err := hist.DoseHistory(ctx, "u1", sameTime.Add(-time.Hour), sameTime.Add(time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(u1AfterEdit))
	for _, d := range u1AfterEdit {
		if d.ID == d1.ID {
			assert.Equal(t, float32(6), d.Dose)
			assert.Equal(t, "edited first", d.Comment)
		} else {
			assert.Equal(t, d2.ID, d.ID)
			assert.Equal(t, float32(5), d.Dose)
			assert.Equal(t, "second", d.Comment)
		}
	}

	// Delete only d2.
	assert.NoError(t, hist.ForgetDoses(ctx, "u1", []uuid.UUID{d2.ID}))
	u1AfterDelete, err := hist.DoseHistory(ctx, "u1", sameTime.Add(-time.Hour), sameTime.Add(time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(u1AfterDelete))
	assert.Equal(t, d1.ID, u1AfterDelete[0].ID)

	// Phase 3: ImportDoses dedup & cross-account restore.
	explicitID := uuid.MustParse("019450a0-0000-7000-8000-000000000111")
	v1Time := time.Date(2025, 3, 1, 8, 0, 0, 0, time.UTC)

	imported, err := hist.ImportDoses(ctx, "u1", slices.Values([]dosage.Dose{
		// Duplicate explicit ID (already exists for u1) -> ignored
		{ID: d1.ID, DeliveryMethod: "EV im", Dose: 6, TakenAt: sameTime, Comment: "edited first"},
		// New explicit ID -> inserted
		{ID: explicitID, DeliveryMethod: "EV im", Dose: 4, TakenAt: v1Time, Comment: "v2 comment"},
		// Duplicate v1 ID-less within same batch -> only 1 inserted
		{DeliveryMethod: "EV im", Dose: 3, TakenAt: v1Time, Comment: "v1 comment"},
		{DeliveryMethod: "EV im", Dose: 3, TakenAt: v1Time, Comment: "v1 comment"},
		// Same timestamp, different comment -> inserted as separate row
		{DeliveryMethod: "EV im", Dose: 3, TakenAt: v1Time, Comment: "v1 different comment"},
	}))
	assert.NoError(t, err)
	assert.Equal(t, int64(3), imported)

	// Re-importing the exact same v1 ID-less row against u1 inserts 0.
	importedAgain, err := hist.ImportDoses(ctx, "u1", slices.Values([]dosage.Dose{
		{DeliveryMethod: "EV im", Dose: 3, TakenAt: v1Time, Comment: "v1 comment"},
	}))
	assert.NoError(t, err)
	assert.Equal(t, int64(0), importedAgain)

	// Cross-account restore: importing d1.ID and explicitID into u2 inserts both.
	importedU2, err := hist.ImportDoses(ctx, "u2", slices.Values([]dosage.Dose{
		{ID: d1.ID, DeliveryMethod: "EV im", Dose: 6, TakenAt: sameTime, Comment: "restored"},
		{ID: explicitID, DeliveryMethod: "EV im", Dose: 4, TakenAt: v1Time, Comment: "v2 comment"},
	}))
	assert.NoError(t, err)
	assert.Equal(t, int64(2), importedU2)

	u2History, err := hist.DoseHistory(ctx, "u2", time.Time{}, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(u2History))
	assert.Equal(t, "restored", u2History[0].Comment)
	assert.Equal(t, "v2 comment", u2History[1].Comment)
}

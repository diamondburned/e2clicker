package postgresql

import (
	"context"
	"fmt"
	"iter"
	"slices"
	"strings"
	"time"
	"uuid"

	"e2clicker.app/internal/sqlc/postgresqlc"
	"e2clicker.app/services/dosage"
	"e2clicker.app/services/user"
	"github.com/jackc/pgx/v5/pgtype"
)

type doseHistoryStorage Storage

func (s *Storage) doseHistoryStorage() dosage.DoseHistoryStorage { return (*doseHistoryStorage)(s) }

func (s *doseHistoryStorage) RecordDose(ctx context.Context, userSecret user.Secret, dose dosage.Dose) (dosage.Dose, error) {
	id := dose.ID
	if id == uuid.Nil() {
		id = uuid.NewV7()
	}

	row, err := s.q.RecordDose(ctx, postgresqlc.RecordDoseParams{
		UserSecret:     userSecret,
		ID:             pgtype.UUID{Bytes: id, Valid: true},
		DeliveryMethod: pgtype.Text{String: dose.DeliveryMethod, Valid: true},
		Dose:           dose.Dose,
		TakenAt:        pgtype.Timestamptz{Time: dose.TakenAt, Valid: true},
		TakenOffAt:     pgtype.Timestamptz{Time: deref(dose.TakenOffAt), Valid: dose.TakenOffAt != nil},
		Comment:        pgtype.Text{String: dose.Comment, Valid: dose.Comment != ""},
	})
	if err != nil {
		return dosage.Dose{}, err
	}
	return convertDose(row), nil
}

func (s *doseHistoryStorage) ImportDoses(ctx context.Context, userSecret user.Secret, doses iter.Seq[dosage.Dose]) (int64, error) {
	const table = "dosage_history"
	insertCols := []string{
		"id",
		"user_secret",
		"delivery_method",
		"dose",
		"taken_at",
		"taken_off_at",
		"comment",
	}
	copyCols := append(slices.Clone(insertCols), "has_explicit_id")

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintln(
		"CREATE TEMP TABLE tmp_history",
		"  (LIKE", table, "INCLUDING DEFAULTS)",
		"ON COMMIT DROP;",
	))
	if err != nil {
		return 0, fmt.Errorf("ImportDoses: create temp table: %w", err)
	}

	if _, err = tx.Exec(ctx, "ALTER TABLE tmp_history ADD COLUMN has_explicit_id boolean NOT NULL DEFAULT false;"); err != nil {
		return 0, fmt.Errorf("ImportDoses: alter temp table: %w", err)
	}

	iter := newCopyFromIterator(doses, func(d dosage.Dose) ([]any, error) {
		hasExplicitID := d.ID != uuid.Nil()
		id := d.ID
		if !hasExplicitID {
			id = uuid.NewV7()
		}

		return []any{
			pgtype.UUID{Bytes: id, Valid: true},
			userSecret,
			d.DeliveryMethod,
			d.Dose,
			pgtype.Timestamptz{Time: d.TakenAt, Valid: true},
			pgtype.Timestamptz{Time: deref(d.TakenOffAt), Valid: d.TakenOffAt != nil},
			pgtype.Text{String: d.Comment, Valid: d.Comment != ""},
			hasExplicitID,
		}, nil
	})
	defer iter.Close()

	if _, err = tx.CopyFrom(ctx, []string{"tmp_history"}, copyCols, iter); err != nil {
		return 0, fmt.Errorf("ImportDoses: copy from: %w", err)
	}

	cols := strings.Join(insertCols, ", ")
	r1, err := tx.Exec(ctx, fmt.Sprintln(
		"INSERT INTO", table, "(", cols, ")",
		"SELECT", cols, "FROM tmp_history",
		"WHERE has_explicit_id",
		"ORDER BY taken_at ASC, id ASC",
		"ON CONFLICT (user_secret, id) DO NOTHING;",
	))
	if err != nil {
		return 0, fmt.Errorf("ImportDoses: insert explicit-id rows: %w", err)
	}

	r2, err := tx.Exec(ctx, fmt.Sprintln(
		"INSERT INTO", table, "(", cols, ")",
		"SELECT DISTINCT ON (taken_at, delivery_method, dose, taken_off_at, comment)",
		cols,
		"FROM tmp_history t",
		"WHERE NOT has_explicit_id",
		"  AND NOT EXISTS (",
		"    SELECT 1 FROM", table, "d",
		"    WHERE d.user_secret = t.user_secret",
		"      AND d.taken_at = t.taken_at",
		"      AND d.delivery_method IS NOT DISTINCT FROM t.delivery_method",
		"      AND d.dose = t.dose",
		"      AND d.taken_off_at IS NOT DISTINCT FROM t.taken_off_at",
		"      AND d.comment IS NOT DISTINCT FROM t.comment",
		"  )",
		"ORDER BY taken_at ASC, delivery_method, dose, taken_off_at, comment, id ASC;",
	))
	if err != nil {
		return 0, fmt.Errorf("ImportDoses: insert v1 rows: %w", err)
	}

	n := r1.RowsAffected() + r2.RowsAffected()

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("ImportDoses: commit: %w", err)
	}

	return n, nil
}

func (s *doseHistoryStorage) EditDose(ctx context.Context, userSecret user.Secret, id uuid.UUID, d dosage.Dose) error {
	n, err := s.q.EditDose(ctx, postgresqlc.EditDoseParams{
		UserSecret:     userSecret,
		ID:             pgtype.UUID{Bytes: id, Valid: true},
		DeliveryMethod: pgtype.Text{String: d.DeliveryMethod, Valid: true},
		Dose:           d.Dose,
		TakenAt:        pgtype.Timestamptz{Time: d.TakenAt, Valid: true},
		TakenOffAt:     pgtype.Timestamptz{Time: deref(d.TakenOffAt), Valid: d.TakenOffAt != nil},
		Comment:        pgtype.Text{String: d.Comment, Valid: d.Comment != ""},
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return dosage.ErrNoDoseMatched
	}
	return nil
}

func (s *doseHistoryStorage) ForgetDoses(ctx context.Context, userSecret user.Secret, ids []uuid.UUID) error {
	pgIDs := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		pgIDs[i] = pgtype.UUID{Bytes: id, Valid: true}
	}

	n, err := s.q.ForgetDoses(ctx, postgresqlc.ForgetDosesParams{
		UserSecret: userSecret,
		IDs:        pgIDs,
	})
	if err != nil {
		return err
	}
	if len(ids) > 0 && n == 0 {
		return dosage.ErrNoDoseMatched
	}
	return nil
}

func (s *doseHistoryStorage) DoseHistory(ctx context.Context, secret user.Secret, begin, end time.Time) ([]dosage.Dose, error) {
	if end.IsZero() {
		end = time.Now()
	}

	rows, err := s.q.DoseHistory(ctx, postgresqlc.DoseHistoryParams{
		UserSecret: secret,
		Start:      pgtype.Timestamptz{Time: begin, Valid: true},
		End:        pgtype.Timestamptz{Time: end, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return convertList(rows, convertDose), nil
}

func convertDose(o postgresqlc.DosageHistory) dosage.Dose {
	return dosage.Dose{
		ID:             uuid.UUID(o.ID.Bytes),
		DeliveryMethod: o.DeliveryMethod.String,
		Dose:           o.Dose,
		TakenAt:        o.TakenAt.Time,
		TakenOffAt:     maybePtr(o.TakenOffAt.Time, o.TakenOffAt.Valid),
		Comment:        o.Comment.String,
	}
}

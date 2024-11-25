package postgresql

import (
	"context"
	"fmt"
	"iter"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"libdb.so/e2clicker/internal/sqlc/postgresqlc"
	"libdb.so/e2clicker/services/dosage"
	"libdb.so/e2clicker/services/user"
)

type doseHistoryStorage Storage

func (s *Storage) doseHistoryStorage() dosage.DoseHistoryStorage { return (*doseHistoryStorage)(s) }

func (s *doseHistoryStorage) RecordDose(ctx context.Context, userSecret user.Secret, dose dosage.Dose) error {
	return s.q.RecordDose(ctx, postgresqlc.RecordDoseParams{
		UserSecret:     userSecret,
		DeliveryMethod: pgtype.Text{String: dose.DeliveryMethod, Valid: true},
		Dose:           dose.Dose,
		TakenAt:        pgtype.Timestamptz{Time: dose.TakenAt, Valid: true},
		TakenOffAt:     pgtype.Timestamptz{Time: deref(dose.TakenOffAt), Valid: dose.TakenOffAt != nil},
	})
}

func (s *doseHistoryStorage) ImportDoses(ctx context.Context, userSecret user.Secret, doses iter.Seq[dosage.Dose]) (int64, error) {
	const table = "dosage_history"
	rows := []string{
		"user_secret",
		"delivery_method",
		"dose",
		"taken_at",
		"taken_off_at",
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintln(
		"CREATE TEMP TABLE tmp_history",
		"  (LIKE", table, "INCLUDING ALL EXCLUDING STORAGE EXCLUDING CONSTRAINTS)",
		"ON COMMIT DROP;",
	))
	if err != nil {
		return 0, fmt.Errorf("ImportDoses: create temp table: %w", err)
	}

	iter := newCopyFromIterator(doses, func(d dosage.Dose) ([]any, error) {
		return []any{
			userSecret,
			d.DeliveryMethod,
			d.Dose,
			pgtype.Timestamptz{Time: d.TakenAt, Valid: true},
			pgtype.Timestamptz{Time: deref(d.TakenOffAt), Valid: d.TakenOffAt != nil},
		}, nil
	})
	defer iter.Close()

	if _, err = tx.CopyFrom(ctx, []string{"tmp_history"}, rows, iter); err != nil {
		return 0, fmt.Errorf("ImportDoses: copy from: %w", err)
	}

	r, err := tx.Exec(ctx, fmt.Sprintln(
		"INSERT INTO", table, "(", strings.Join(rows, ", "), ")",
		"SELECT", strings.Join(rows, ", "), "FROM tmp_history ORDER BY taken_at ASC",
		"ON CONFLICT (user_secret, taken_at) DO NOTHING;",
	))

	n := r.RowsAffected()

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("ImportDoses: commit: %w", err)
	}

	return n, nil
}

func (s *doseHistoryStorage) EditDose(ctx context.Context, userSecret user.Secret, doseTime time.Time, d dosage.Dose) error {
	n, err := s.q.EditDose(ctx, postgresqlc.EditDoseParams{
		UserSecret: userSecret,
		OldTakenAt: pgtype.Timestamptz{Time: doseTime, Valid: true},

		DeliveryMethod: pgtype.Text{String: d.DeliveryMethod, Valid: true},
		Dose:           d.Dose,
		TakenAt:        pgtype.Timestamptz{Time: d.TakenAt, Valid: true},
		TakenOffAt:     pgtype.Timestamptz{Time: deref(d.TakenOffAt), Valid: d.TakenOffAt != nil},
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return dosage.ErrNoDoseMatched
	}
	return nil
}

func (s *doseHistoryStorage) ForgetDoses(ctx context.Context, userSecret user.Secret, doseTimes []time.Time) error {
	pgTimes := make([]pgtype.Timestamp, len(doseTimes))
	for i, t := range doseTimes {
		pgTimes[i] = pgtype.Timestamp{Time: t, Valid: true}
	}

	n, err := s.q.ForgetDoses(ctx, postgresqlc.ForgetDosesParams{
		UserSecret: userSecret,
		TakenAt:    pgTimes,
	})
	if err != nil {
		return err
	}
	if len(doseTimes) > 0 && n == 0 {
		return dosage.ErrNoDoseMatched
	}
	return nil
}

func (s *doseHistoryStorage) DoseHistory(ctx context.Context, secret user.Secret, begin, end time.Time) iter.Seq2[dosage.Dose, error] {
	if end.IsZero() {
		end = time.Now()
	}

	iter := s.q.DoseHistory(ctx, postgresqlc.DoseHistoryParams{
		UserSecret: secret,
		Start:      pgtype.Timestamptz{Time: begin, Valid: true},
		End:        pgtype.Timestamptz{Time: end, Valid: true},
	})

	return func(yield func(dosage.Dose, error) bool) {
		for o1 := range iter.Iterate() {
			o2 := convertDose(o1)
			if !yield(o2, nil) {
				return
			}
		}

		if err := iter.Err(); err != nil {
			yield(dosage.Dose{}, err)
		}
	}
}

func convertDose(o postgresqlc.DosageHistory) dosage.Dose {
	return dosage.Dose{
		DeliveryMethod: o.DeliveryMethod.String,
		Dose:           o.Dose,
		TakenAt:        o.TakenAt.Time,
		TakenOffAt:     maybePtr(o.TakenOffAt.Time, o.TakenOffAt.Valid),
		Comment:        o.Comment.String,
	}
}

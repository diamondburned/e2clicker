package postgresql

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"libdb.so/e2clicker/internal/sqlc"
	"libdb.so/e2clicker/internal/sqlc/postgresqlc"
	"libdb.so/e2clicker/services/dosage"
	"libdb.so/e2clicker/services/user"
)

func (s *Storage) dosageStorage() dosage.DosageStorage { return (*dosageStorage)(s) }

type dosageStorage Storage

func (s *dosageStorage) DeliveryMethods(ctx context.Context) ([]dosage.DeliveryMethod, error) {
	methods, err := s.q.DeliveryMethods(ctx)
	if err != nil {
		return nil, err
	}

	return convertList(methods, func(m postgresqlc.DeliveryMethod) dosage.DeliveryMethod {
		return dosage.DeliveryMethod(m)
	}), nil
}

func (s *dosageStorage) DosageSchedule(ctx context.Context, secret user.Secret) (*dosage.Schedule, error) {
	d, err := s.q.DosageSchedule(ctx, sqlc.XID(secret))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var interval dosage.Days = 0 +
		(dosage.Days(d.Interval.Days)) +
		(dosage.Days(d.Interval.Microseconds) / 1e6 / (60 * 60 * 24)) +
		(dosage.Days(d.Interval.Months) * 30)

	return &dosage.Schedule{
		UserSecret:     secret,
		DeliveryMethod: d.DeliveryMethod.String,
		Dose:           d.Dose,
		Interval:       interval,
		Concurrence:    maybePtr(int(d.Concurrence.Int16), d.Concurrence.Valid),
	}, nil
}

func (s *dosageStorage) SetDosageSchedule(ctx context.Context, d dosage.Schedule) error {
	int, frac := math.Modf(float64(d.Interval))
	return s.q.SetDosageSchedule(ctx, postgresqlc.SetDosageScheduleParams{
		UserSecret:     sqlc.XID(d.UserSecret),
		DeliveryMethod: pgtype.Text{String: d.DeliveryMethod, Valid: true},
		Dose:           d.Dose,
		Interval: pgtype.Interval{
			Days:         int32(int),
			Microseconds: int64(frac * 24 * 60 * 60 * 1e6),
			Valid:        true,
		},
		Concurrence: pgtype.Int2{
			Int16: int16(min(deref(d.Concurrence), math.MaxInt16)),
			Valid: d.Concurrence != nil && *d.Concurrence > 0,
		},
	})
}

func (s *dosageStorage) ClearDosageSchedule(ctx context.Context, secret user.Secret) error {
	return s.q.DeleteDosageSchedule(ctx, sqlc.XID(secret))
}

func (s *dosageStorage) RecordDose(ctx context.Context, userSecret user.Secret, takenAt time.Time) (dosage.Observation, error) {
	o, err := s.q.RecordDose(ctx, postgresqlc.RecordDoseParams{
		UserSecret: sqlc.XID(userSecret),
		TakenAt:    pgtype.Timestamptz{Time: takenAt, Valid: true},
	})
	if err != nil {
		return dosage.Observation{}, err
	}
	return convertObservation(o), nil
}

func (s *dosageStorage) EditDose(ctx context.Context, userSecret user.Secret, o dosage.Observation) error {
	n, err := s.q.EditDose(ctx, postgresqlc.EditDoseParams{
		UserSecret:     sqlc.XID(userSecret),
		DoseID:         o.DoseID,
		DeliveryMethod: pgtype.Text{String: o.DeliveryMethod, Valid: true},
		Dose:           o.Dose,
		TakenAt:        pgtype.Timestamptz{Time: o.TakenAt, Valid: true},
		TakenOffAt:     pgtype.Timestamptz{Time: deref(o.TakenOffAt), Valid: o.TakenOffAt != nil},
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return dosage.ErrNoDoseMatched
	}
	return nil
}

func (s *dosageStorage) ForgetDoses(ctx context.Context, userSecret user.Secret, doseIDs []int64) error {
	n, err := s.q.ForgetDoses(ctx, postgresqlc.ForgetDosesParams{
		UserSecret: sqlc.XID(userSecret),
		DoseIDs:    doseIDs,
	})
	if err != nil {
		return err
	}
	if len(doseIDs) > 0 && n == 0 {
		return dosage.ErrNoDoseMatched
	}
	return nil
}

func (s *dosageStorage) DoseHistory(ctx context.Context, secret user.Secret, begin, end time.Time) (dosage.History, error) {
	o, err := s.q.DoseHistory(ctx, postgresqlc.DoseHistoryParams{
		UserSecret: sqlc.XID(secret),
		Start:      pgtype.Timestamptz{Time: begin, Valid: true},
		End:        pgtype.Timestamptz{Time: end, Valid: true},
	})
	if err != nil {
		return dosage.History{}, err
	}
	return dosage.History{
		UserSecret: secret,
		Entries:    convertList(o, convertObservation),
	}, nil
}

func convertObservation(o postgresqlc.DosageHistory) dosage.Observation {
	return dosage.Observation{
		DoseID:         o.DoseID,
		DeliveryMethod: o.DeliveryMethod.String,
		Dose:           o.Dose,
		TakenAt:        o.TakenAt.Time,
		TakenOffAt:     maybePtr(o.TakenOffAt.Time, o.TakenOffAt.Valid),
	}
}

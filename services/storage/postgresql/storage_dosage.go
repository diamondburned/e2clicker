package postgresql

import (
	"context"
	"errors"
	"math"

	"e2clicker.app/internal/sqlc/postgresqlc"
	"e2clicker.app/services/dosage"
	"e2clicker.app/services/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type dosageStorage Storage

func (s *Storage) dosageStorage() dosage.DosageStorage { return (*dosageStorage)(s) }

func (s *dosageStorage) DeliveryMethods(ctx context.Context) ([]dosage.DeliveryMethod, error) {
	methods, err := s.q.DeliveryMethods(ctx)
	if err != nil {
		return nil, err
	}

	return convertList(methods, func(m postgresqlc.DeliveryMethod) dosage.DeliveryMethod {
		return dosage.DeliveryMethod(m)
	}), nil
}

func (s *dosageStorage) Dosage(ctx context.Context, secret user.Secret) (*dosage.Dosage, error) {
	d, err := s.q.DosageSchedule(ctx, secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	d2 := convertDosage(d)
	return &d2, nil
}

func convertDosage(d postgresqlc.DosageSchedule) dosage.Dosage {
	interval := dosage.Days(0) +
		(dosage.Days(d.Interval.Days)) +
		(dosage.Days(d.Interval.Microseconds) / 1e6 / (60 * 60 * 24)) +
		(dosage.Days(d.Interval.Months) * 30)

	return dosage.Dosage{
		UserSecret:     d.UserSecret,
		DeliveryMethod: d.DeliveryMethod.String,
		Dose:           d.Dose,
		Interval:       interval,
		Concurrence:    maybePtr(int(d.Concurrence.Int16), d.Concurrence.Valid),
	}
}

func (s *dosageStorage) SetDosage(ctx context.Context, d dosage.Dosage) error {
	int, frac := math.Modf(float64(d.Interval))
	return s.q.SetDosageSchedule(ctx, postgresqlc.SetDosageScheduleParams{
		UserSecret:     d.UserSecret,
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

func (s *dosageStorage) ClearDosage(ctx context.Context, secret user.Secret) error {
	return s.q.DeleteDosageSchedule(ctx, secret)
}

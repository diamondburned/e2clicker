package postgresql

import (
	"context"
	"errors"
	"math"

	"e2clicker.app/internal/meta"
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

	return &dosage.Dosage{
		DeliveryMethod:     d.DeliveryMethod.String,
		Dose:               d.Dose,
		Interval:           meta.DaysFromPostgreSQLInterval(d.Interval),
		Concurrence:        maybePtr(int(d.Concurrence.Int16), d.Concurrence.Valid),
		ReminderRecurrence: convertList(d.ReminderRecurrence, meta.DaysFromPostgreSQLInterval),
	}, nil
}

func (s *dosageStorage) SetDosage(ctx context.Context, secret user.Secret, d dosage.Dosage) error {
	return s.q.SetDosageSchedule(ctx, postgresqlc.SetDosageScheduleParams{
		UserSecret:     secret,
		DeliveryMethod: pgtype.Text{String: d.DeliveryMethod, Valid: true},
		Dose:           d.Dose,
		Interval:       d.Interval.ToPostgreSQLInterval(),
		Concurrence: pgtype.Int2{
			Int16: int16(min(deref(d.Concurrence), math.MaxInt16)),
			Valid: d.Concurrence != nil && *d.Concurrence > 0,
		},
		ReminderRecurrence: convertList(d.ReminderRecurrence, meta.Days.ToPostgreSQLInterval),
	})
}

func (s *dosageStorage) ClearDosage(ctx context.Context, secret user.Secret) error {
	return s.q.DeleteDosageSchedule(ctx, secret)
}

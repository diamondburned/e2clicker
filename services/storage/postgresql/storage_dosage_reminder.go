package postgresql

import (
	"context"
	"errors"

	"e2clicker.app/internal/meta"
	"e2clicker.app/internal/ptr"
	"e2clicker.app/internal/sqlc/postgresqlc"
	"e2clicker.app/services/dosage"
	"github.com/jackc/pgx/v5/pgtype"
)

type dosageReminderStorage Storage

func (s *Storage) dosageReminderStorage() dosage.DosageReminderStorage {
	return (*dosageReminderStorage)(s)
}

func (s *dosageReminderStorage) UpcomingDosageReminders(ctx context.Context) ([]dosage.DosageReminder, error) {
	rows, err := s.q.UpcomingDosageReminders(ctx)
	if err != nil {
		return nil, err
	}

	return convertList(rows, func(m postgresqlc.UpcomingDosageRemindersRow) dosage.DosageReminder {
		return dosage.DosageReminder{
			UserSecret: m.UserSecret,
			Username:   m.UserName,
			Dosage: dosage.Dosage{
				DeliveryMethod:     m.DosageSchedule.DeliveryMethod.String,
				Dose:               m.DosageSchedule.Dose,
				Interval:           meta.DaysFromPostgreSQLInterval(m.DosageSchedule.Interval),
				Concurrence:        maybePtr(int(m.DosageSchedule.Concurrence.Int16), m.DosageSchedule.Concurrence.Valid),
				ReminderRecurrence: convertList(m.DosageSchedule.ReminderRecurrence, meta.DaysFromPostgreSQLInterval),
			},
			LastDose:         convertDose(m.DosageHistory),
			LastRemindedDose: ptr.ToIf(m.LastRemindedDose.Time, m.LastRemindedDose.Valid),
			LastRemindedAt:   ptr.ToIf(m.LastRemindedAt.Time, m.LastRemindedAt.Valid),
			SnoozedUntil:     nil, // TODO: support snoozed until
		}
	}), nil
}

func (s *dosageReminderStorage) RecordRemindedDoseAttempts(ctx context.Context, remindedDoseAttempts []dosage.RemindedDoseAttempt) error {
	var errs []error
	for _, attempt := range remindedDoseAttempts {
		err := s.q.RecordRemindedDoseAttempt(ctx, postgresqlc.RecordRemindedDoseAttemptParams{
			UserSecret:         attempt.UserSecret,
			SentAt:             pgtype.Timestamptz{Time: attempt.RemindedAt, Valid: true},
			SupposedEntityTime: pgtype.Timestamptz{Time: attempt.RemindedDose, Valid: true},
			ErrorReason:        pgtype.Text{String: ptr.Deref(attempt.ErrorReason), Valid: attempt.ErrorReason != nil},
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

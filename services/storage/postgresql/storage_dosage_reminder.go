package postgresql

import (
	"context"
	"errors"
	"iter"

	"e2clicker.app/internal/ptr"
	"e2clicker.app/internal/sqlc/postgresqlc"
	"e2clicker.app/services/dosage"
	"github.com/jackc/pgx/v5/pgtype"
)

type dosageReminderStorage Storage

func (s *Storage) dosageReminderStorage() dosage.DosageReminderStorage {
	return (*dosageReminderStorage)(s)
}

func (s *dosageReminderStorage) UpcomingDosageReminders(ctx context.Context) iter.Seq2[dosage.DosageReminder, error] {
	iter := s.q.UpcomingDosageReminders(ctx)

	return func(yield func(dosage.DosageReminder, error) bool) {
		for o1 := range iter.Iterate() {
			o2 := dosage.DosageReminder{
				UserSecret:       o1.UserSecret,
				Username:         o1.UserName,
				Dosage:           convertDosage(o1.DosageSchedule),
				LastDose:         convertDose(o1.DosageHistory),
				LastRemindedDose: ptr.ToIf(o1.LastNotificationTime.Time, o1.LastNotificationTime.Valid),
				SnoozedUntil:     nil,
			}

			if !yield(o2, nil) {
				return
			}
		}

		if err := iter.Err(); err != nil {
			yield(dosage.DosageReminder{}, err)
		}
	}
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

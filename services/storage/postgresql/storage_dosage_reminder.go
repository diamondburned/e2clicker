package postgresql

import (
	"context"
	"iter"

	"e2clicker.app/services/dosage"
)

type dosageReminderStorage Storage

func (s *Storage) dosageReminderStorage() dosage.DosageReminderStorage {
	return (*dosageReminderStorage)(s)
}

func (s *dosageReminderStorage) UpcomingDosageReminders(ctx context.Context) iter.Seq2[dosage.DosageReminder, error] {
	panic("TODO")
}

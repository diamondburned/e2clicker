package postgresql

import (
	"context"
	"testing"
	"time"

	"e2clicker.app/internal/meta"
	"e2clicker.app/internal/ptr"
	e2clickermodule "e2clicker.app/nix/modules/e2clicker"
	"e2clicker.app/services/dosage"
	"github.com/alecthomas/assert/v2"
	"github.com/neilotoole/slogt"
	"go.uber.org/fx/fxtest"
)

func TestNotificationAttemptAndUpcomingReminders(t *testing.T) {
	ctx := context.Background()
	uri := startTestPostgres(t)

	lc := fxtest.NewLifecycle(t)
	s, err := NewStorage(lc, e2clickermodule.PostgreSQL{DatabaseURI: uri}, slogt.New(t))
	assert.NoError(t, err)
	lc.RequireStart()
	defer lc.RequireStop()

	const secret = "test-user-1"
	_, err = s.userStorage().CreateUser(ctx, secret, "Test User")
	assert.NoError(t, err)

	notifStorage := s.notificationUserStorage()
	dosageStorage := s.dosageStorage()
	doseHistStorage := s.doseHistoryStorage()
	reminderStorage := s.dosageReminderStorage()

	// No notification attempts yet.
	attempt, err := notifStorage.LastNotificationAttempt(ctx, secret)
	assert.NoError(t, err)
	assert.Zero(t, attempt)

	takenAt := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	assert.NoError(t, dosageStorage.SetDosage(ctx, secret, dosage.Dosage{
		DeliveryMethod: "EV im",
		Dose:           5,
		Interval:       meta.Days(1),
	}))
	_, err = doseHistStorage.RecordDose(ctx, secret, dosage.Dose{
		DeliveryMethod: "EV im",
		Dose:           5,
		TakenAt:        takenAt,
	})
	assert.NoError(t, err)

	// Record an errored reminder attempt; it must be counted by UpcomingDosageReminders
	// and exposed by LastNotificationAttempt.
	firstAttemptAt := takenAt.Add(24 * time.Hour)
	errReason := "webPush: push subscription expired"
	assert.NoError(t, reminderStorage.RecordRemindedDoseAttempts(ctx, []dosage.RemindedDoseAttempt{
		{
			UserSecret:   secret,
			RemindedAt:   firstAttemptAt,
			RemindedDose: takenAt,
			ErrorReason:  ptr.To(errReason),
		},
	}))

	attempt, err = notifStorage.LastNotificationAttempt(ctx, secret)
	assert.NoError(t, err)
	assert.True(t, attempt != nil)
	assert.True(t, firstAttemptAt.Equal(attempt.SentAt))
	assert.True(t, attempt.Errored)
	assert.Equal(t, ptr.To(errReason), attempt.ErrorReason)

	upcoming, err := reminderStorage.UpcomingDosageReminders(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(upcoming))
	assert.NotZero(t, upcoming[0].LastRemindedAt)
	assert.True(t, firstAttemptAt.Equal(*upcoming[0].LastRemindedAt))

	// A newer successful attempt updates LastNotificationAttempt to non-errored.
	secondAttemptAt := firstAttemptAt.Add(time.Hour)
	assert.NoError(t, reminderStorage.RecordRemindedDoseAttempts(ctx, []dosage.RemindedDoseAttempt{
		{
			UserSecret:   secret,
			RemindedAt:   secondAttemptAt,
			RemindedDose: takenAt,
		},
	}))

	attempt, err = notifStorage.LastNotificationAttempt(ctx, secret)
	assert.NoError(t, err)
	assert.True(t, attempt != nil)
	assert.True(t, secondAttemptAt.Equal(attempt.SentAt))
	assert.False(t, attempt.Errored)
	assert.Zero(t, attempt.ErrorReason)
}

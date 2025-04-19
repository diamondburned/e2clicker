package dosage

import (
	"context"
	"fmt"
	"iter"
	"testing"
	"time"

	"e2clicker.app/internal/meta"
	"e2clicker.app/internal/ptr"
	"e2clicker.app/services/notification/openapi"
	"e2clicker.app/services/user"
	"github.com/alecthomas/assert/v2"
	"github.com/neilotoole/slogt"
	"go.uber.org/fx/fxtest"
)

//go:generate moq -out reminder_mock_test.go . DosageReminderStorage
//go:generate moq -out reminder_mock_notification_test.go -pkg dosage ../notification UserNotificationService

const day = 24 * time.Hour

var now = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

func TestIngestReminders(t *testing.T) {
	now := now

	type userSet map[string]struct{}

	newUserSet := func(users ...string) userSet {
		result := make(userSet, len(users))
		for _, u := range users {
			result[u] = struct{}{}
		}
		return result
	}

	testCases := []struct {
		name            string
		checkTime       time.Time
		reminders       []DosageReminder
		remindedUsers   userSet
		expectedNextRun time.Time
	}{
		{
			name:            "empty",
			reminders:       []DosageReminder{},
			remindedUsers:   newUserSet(),
			expectedNextRun: now.Add(nextUpdateInterval),
		},
		{
			name: "not_relevant",
			reminders: []DosageReminder{
				{
					Username: "user1",
					Dosage:   Dosage{Interval: 1},
					LastDose: Dose{TakenAt: now},
				},
			},
			remindedUsers:   newUserSet(),
			expectedNextRun: now.Add(nextUpdateInterval),
		},
		{
			name:      "one_relevant_nearest",
			checkTime: now.Add(-1 * time.Minute), // 1 minute before the next dose
			reminders: []DosageReminder{
				{
					Username: "user1",
					Dosage:   Dosage{Interval: 1},
					LastDose: Dose{TakenAt: now.Add(-day)},
				},
			},
			remindedUsers: newUserSet(),
			// expectedNextRun: now.Add(shortestNextNotification),
		},
		{
			name:      "one_relevant_near_enough",
			checkTime: now.Add(-10 * time.Minute), // 10 minutes before the next dose
			reminders: []DosageReminder{
				{
					Username: "user1",
					Dosage:   Dosage{Interval: 1},
					LastDose: Dose{TakenAt: now.Add(-day)},
				},
			},
			remindedUsers: newUserSet(),
			// expectedNextRun: now.Add(10 * time.Minute),
		},
		{
			name: "one_relevant_too_early",
			reminders: []DosageReminder{
				{
					Username: "user1",
					Dosage:   Dosage{Interval: 1},
					LastDose: Dose{TakenAt: now.Add(-day + nextUpdateInterval + time.Minute)}, // 30 minutes before dose
				},
			},
			remindedUsers: newUserSet(),
			// expectedNextRun: now.Add(nextUpdateInterval),
		},
		{
			name: "one_notifying",
			reminders: []DosageReminder{
				{
					Username: "user1",
					Dosage:   Dosage{Interval: 1},
					LastDose: Dose{TakenAt: now.Add(-day - time.Minute)}, // 1 minute after dose
				},
			},
			remindedUsers: newUserSet("user1"),
			// expectedNextRun: now.Add(nextUpdateInterval),
		},
		{
			name: "one_notifying_but_already_notified",
			reminders: []DosageReminder{
				{
					Username:       "user1",
					Dosage:         Dosage{Interval: 1},
					LastDose:       Dose{TakenAt: now.Add(-day - time.Minute)}, // 1 minute after dose
					LastRemindedAt: ptr.To(now.Add(-time.Minute)),
				},
			},
			remindedUsers: newUserSet(),
			// expectedNextRun: now.Add(nextUpdateInterval),
		},
		{
			name: "one_notifying_but_snoozed",
			reminders: []DosageReminder{
				{
					Username:     "user1",
					Dosage:       Dosage{Interval: 1},
					LastDose:     Dose{TakenAt: now.Add(-day - time.Minute)}, // 1 minute after dose
					SnoozedUntil: ptr.To(now.Add(10 * time.Minute)),
				},
			},
			remindedUsers: newUserSet(),
			// expectedNextRun: now.Add(10 * time.Minute),
		},
		{
			name:      "one_notifying_recurrent_initial",
			checkTime: now.Add(1*day + 1),
			reminders: []DosageReminder{
				{
					Username: "user1",
					Dosage:   Dosage{Interval: 1, ReminderRecurrence: []meta.Days{1, 2}},
					LastDose: Dose{TakenAt: now},
				},
			},
			remindedUsers: newUserSet("user1"),
		},
		{
			name:      "one_notifying_recurrent_initial_done",
			checkTime: now.Add(1*day + 100),
			reminders: []DosageReminder{
				{
					Username:       "user1",
					Dosage:         Dosage{Interval: 1, ReminderRecurrence: []meta.Days{1, 2}},
					LastDose:       Dose{TakenAt: now},
					LastRemindedAt: ptr.To(now.Add(1*day + 1)), // previous checkTime
				},
			},
			remindedUsers: newUserSet(),
		},
		{
			name:      "one_notifying_recurrent_first",
			checkTime: now.Add(2*day + 1),
			reminders: []DosageReminder{
				{
					Username:       "user1",
					Dosage:         Dosage{Interval: 1, ReminderRecurrence: []meta.Days{1, 2}},
					LastDose:       Dose{TakenAt: now},
					LastRemindedAt: ptr.To(now.Add(1*day + 1)),
				},
			},
			remindedUsers: newUserSet("user1"),
		},
		{
			name:      "one_notifying_recurrent_first_done",
			checkTime: now.Add(2*day + 100),
			reminders: []DosageReminder{
				{
					Username:       "user1",
					Dosage:         Dosage{Interval: 1, ReminderRecurrence: []meta.Days{1, 2}},
					LastDose:       Dose{TakenAt: now},
					LastRemindedAt: ptr.To(now.Add(2*day + 1)),
				},
			},
			remindedUsers: newUserSet(),
		},
		{
			name:      "one_notifying_recurrent_second",
			checkTime: now.Add(3*day + 1),
			reminders: []DosageReminder{
				{
					Username:       "user1",
					Dosage:         Dosage{Interval: 1, ReminderRecurrence: []meta.Days{1, 2}},
					LastDose:       Dose{TakenAt: now},
					LastRemindedAt: ptr.To(now.Add(2*day + 1)),
				},
			},
			remindedUsers: newUserSet("user1"),
		},
		{
			name:      "one_notifying_recurrent_second_done",
			checkTime: now.Add(3*day + 100),
			reminders: []DosageReminder{
				{
					Username:       "user1",
					Dosage:         Dosage{Interval: 1, ReminderRecurrence: []meta.Days{1, 2}},
					LastDose:       Dose{TakenAt: now},
					LastRemindedAt: ptr.To(now.Add(3*day + 1)),
				},
			},
			remindedUsers: newUserSet(),
		},
		{
			name:      "one_notifying_recurrent_second_long_after",
			checkTime: now.Add(10 * day),
			reminders: []DosageReminder{
				{
					Username:       "user1",
					Dosage:         Dosage{Interval: 1, ReminderRecurrence: []meta.Days{1, 2}},
					LastDose:       Dose{TakenAt: now},
					LastRemindedAt: ptr.To(now.Add(3*day + 1)),
				},
			},
			remindedUsers: newUserSet(),
		},
	}

	for _, tc := range testCases {
		if tc.checkTime.IsZero() {
			tc.checkTime = now
		}

		t.Run(tc.name, func(t *testing.T) {
			remindersIter := func(yield func(DosageReminder, error) bool) {
				for _, r := range tc.reminders {
					if !yield(r, nil) {
						return
					}
				}
			}

			tracked, err := ingestReminders(tc.checkTime, remindersIter, slogt.New(t))
			assert.NoError(t, err, "ingestReminders must not return an error")

			relevantUsers := make(userSet, len(tracked.notifyingReminders))
			for _, r := range tracked.notifyingReminders {
				relevantUsers[r.Username] = struct{}{}
			}
			assert.Equal(t, tc.remindedUsers, relevantUsers, "ingestReminders must return the expected result")

			if !tc.expectedNextRun.IsZero() {
				assert.Equal(t, tc.expectedNextRun, tracked.nextRun, "ingestReminders must return the expected next run")
			}
		})
	}
}

func TestDosageReminderService(t *testing.T) {
	t.Run("assert_no_duplicate", func(t *testing.T) {
		now := now
		ctx := context.Background()

		s := newMockDosageReminderService(t)

		notified := map[user.Secret]int{}
		s.UserNotificationService.NotifyUserFunc = func(ctx context.Context, secret user.Secret, t openapi.NotificationType) error {
			notified[secret]++
			return nil
		}

		now = s.cycle(ctx, now)
		now = s.cycle(ctx, now)

		assert.Equal(t, 1, len(notified), "only one user should be notified")
		assert.Equal(t, 1, notified["user1"], "user1 should be notified once")
	})
}

type mockDosageReminderService struct {
	*DosageReminderService
	UserNotificationService *UserNotificationServiceMock
	DosageReminderStorage   *mockDosageReminderStorage
}

func newMockDosageReminderService(t *testing.T) *mockDosageReminderService {
	slog := slogt.New(t)

	storage := newMockDosageReminderStorage([]DosageReminder{
		{
			UserSecret: "user1",
			Dosage:     Dosage{Interval: 1},
			LastDose:   Dose{TakenAt: now.Add(-day)},
		},
	})

	notifService := &UserNotificationServiceMock{}

	lc := fxtest.NewLifecycle(t)
	return &mockDosageReminderService{
		DosageReminderService:   NewDosageReminderService(storage, notifService, slog, lc),
		UserNotificationService: notifService,
		DosageReminderStorage:   storage,
	}
}

type mockDosageReminderStorage struct {
	upcoming map[user.Secret]DosageReminder
}

func newMockDosageReminderStorage(upcoming []DosageReminder) *mockDosageReminderStorage {
	s := make(map[user.Secret]DosageReminder, len(upcoming))
	for _, r := range upcoming {
		s[r.UserSecret] = r
	}
	return &mockDosageReminderStorage{upcoming: s}
}

func (m *mockDosageReminderStorage) UpcomingDosageReminders(ctx context.Context) iter.Seq2[DosageReminder, error] {
	return func(yield func(DosageReminder, error) bool) {
		for _, r := range m.upcoming {
			if !yield(r, nil) {
				return
			}
		}
	}
}

func (m *mockDosageReminderStorage) RecordRemindedDoseAttempts(ctx context.Context, attempts []RemindedDoseAttempt) error {
	for _, attempt := range attempts {
		reminder, ok := m.upcoming[attempt.UserSecret]
		if !ok {
			return fmt.Errorf("reminder not found for user %s", attempt.UserSecret)
		}

		reminder.LastRemindedAt = ptr.To(attempt.RemindedAt)
		reminder.LastRemindedDose = ptr.To(attempt.RemindedDose)
		reminder.LastDose.TakenAt = attempt.RemindedDose

		m.upcoming[attempt.UserSecret] = reminder
	}

	return nil
}

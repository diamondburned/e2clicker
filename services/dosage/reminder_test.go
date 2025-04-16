package dosage

import (
	"testing"
	"time"

	"e2clicker.app/internal/meta"
	"e2clicker.app/internal/ptr"
	"github.com/alecthomas/assert/v2"
	"github.com/neilotoole/slogt"
)

func TestIngestReminders(t *testing.T) {
	const day = 24 * time.Hour

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

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
			name: "one_relevant_nearest",
			reminders: []DosageReminder{
				{
					Username: "user1",
					Dosage:   Dosage{Interval: 1},
					LastDose: Dose{TakenAt: now.Add(-day + time.Minute)}, // 1 minute before dose
				},
			},
			remindedUsers: newUserSet(),
			// expectedNextRun: now.Add(shortestNextNotification),
		},
		{
			name: "one_relevant_near_enough",
			reminders: []DosageReminder{
				{
					Username: "user1",
					Dosage:   Dosage{Interval: 1},
					LastDose: Dose{TakenAt: now.Add(-day + 10*time.Minute)}, // 10 minutes before dose
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
					LastRemindedAt: ptr.To(now.Add(0*day + 1)),
				},
			},
			remindedUsers: newUserSet("user1"),
		},
		{
			name:      "one_notifying_recurrent_first",
			checkTime: now.Add(2*day + 1),
			reminders: []DosageReminder{
				{
					Username:       "user1",
					Dosage:         Dosage{Interval: 1, ReminderRecurrence: []meta.Days{1, 2}},
					LastDose:       Dose{TakenAt: now},
					LastRemindedAt: ptr.To(now.Add(1*day + 0)),
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
					LastRemindedAt: ptr.To(now.Add(2*day + 0)),
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

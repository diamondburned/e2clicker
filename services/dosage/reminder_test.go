package dosage

import (
	"testing"
	"time"

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
			remindedUsers:   newUserSet(),
			expectedNextRun: now.Add(shortestNextNotification),
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
			remindedUsers:   newUserSet(),
			expectedNextRun: now.Add(10 * time.Minute),
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
			remindedUsers:   newUserSet(),
			expectedNextRun: now.Add(nextUpdateInterval),
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
			remindedUsers:   newUserSet("user1"),
			expectedNextRun: now.Add(nextUpdateInterval),
		},
		{
			name: "one_notifying_but_already_notified",
			reminders: []DosageReminder{
				{
					Username:         "user1",
					Dosage:           Dosage{Interval: 1},
					LastDose:         Dose{TakenAt: now.Add(-day - time.Minute)}, // 1 minute after dose
					LastRemindedDose: ptr.To(now.Add(-day - time.Minute)),
				},
			},
			remindedUsers:   newUserSet(),
			expectedNextRun: now.Add(nextUpdateInterval),
		},
		{
			name: "one_notifying_but_snoozed",
			reminders: []DosageReminder{
				{
					Username:         "user1",
					Dosage:           Dosage{Interval: 1},
					LastDose:         Dose{TakenAt: now.Add(-day - time.Minute)}, // 1 minute after dose
					LastRemindedDose: ptr.To(now.Add(-day - time.Minute)),
					SnoozedUntil:     ptr.To(now.Add(10 * time.Minute)),
				},
			},
			remindedUsers:   newUserSet(),
			expectedNextRun: now.Add(10 * time.Minute),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			remindersIter := func(yield func(DosageReminder, error) bool) {
				for _, r := range tc.reminders {
					if !yield(r, nil) {
						return
					}
				}
			}

			tracked, err := ingestReminders(now, remindersIter, slogt.New(t))
			assert.NoError(t, err, "ingestReminders must not return an error")

			relevantUsers := make(userSet, len(tracked.notifyingReminders))
			for _, r := range tracked.notifyingReminders {
				relevantUsers[r.Username] = struct{}{}
			}
			assert.Equal(t, tc.remindedUsers, relevantUsers, "ingestReminders must return the expected result")
			assert.Equal(t, tc.expectedNextRun, tracked.nextRun, "ingestReminders must return the expected next run")
		})
	}
}

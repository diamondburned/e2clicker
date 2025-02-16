package dosage

import (
	"context"
	"iter"
	"log/slog"
	"time"

	"e2clicker.app/services/notification"
	notificationapi "e2clicker.app/services/notification/openapi"
	"e2clicker.app/services/user"
	"go.uber.org/fx"
)

// DosageReminderStorage is a storage for dosage reminder data.
type DosageReminderStorage interface {
	// UpcomingDosageReminders returns the upcoming dosage reminders as a
	// streaming iterator. A satisfying but least optimized implementation
	// can return every single person with a dosage schedule and history.
	//
	// No strict ordering is required, but the reminders should be returned in
	// the order of the last dosage time.
	//
	// Users with no dosage schedule or history should not be included in the
	// results.
	UpcomingDosageReminders(ctx context.Context) iter.Seq2[DosageReminder, error]

	// RecordRemindedDoses records the reminded doses.
	// This is used to mark the reminder as sent.
	RecordRemindedDoses(ctx context.Context, remindedDoses []RemindedDose) error
}

// DosageReminder is a reminder for a dosage.
type DosageReminder struct {
	// UserSecret is the secret of the user.
	UserSecret user.Secret
	// Username is the username of the user.
	Username string
	// Dosage is the dosage information of the user.
	Dosage Dosage
	// LastDose is the last dose taken by the user.
	LastDose Dose
	// LastRemindedDose is the TakenAt time of the last reminded dose.
	// This field is optional and is only set if the reminder was recorded.
	LastRemindedDose *time.Time
	// SnoozedUntil is the time until the reminder is snoozed.
	// This field is optional and is only set if the reminder is snoozed by the
	// user on the last notification.
	SnoozedUntil *time.Time
}

// RemindedDose is a dose that has been reminded.
type RemindedDose struct {
	// UserSecret is the secret of the user.
	UserSecret user.Secret
	// RemindedDose is the dose that was reminded.
	// This is the TakenAt time of the dose.
	RemindedDose time.Time
	// ClearSnooze is true if the snooze should be cleared.
	// This is the case if the reminder was sent at the snoozed time.
	ClearSnooze bool
}

// NextNotification returns the time of the next notification.
// It returns the snoozed time if the reminder is snoozed, otherwise it returns
// the time of the last dose plus the interval.
func (r DosageReminder) NextNotification() time.Time {
	if r.SnoozedUntil != nil {
		return *r.SnoozedUntil
	}
	return r.LastDose.TakenAt.Add(r.Dosage.Interval.ToDuration())
}

const (
	shortestNextNotification  = 5 * time.Minute
	nextUpdateInterval        = 30 * time.Minute
	nextUpdateIntervalOnError = 2 * time.Minute
)

// DosageReminderService is a service for managing dosage reminders.
type DosageReminderService struct {
	storage DosageReminderStorage
	notifs  *notification.UserNotificationService
}

// NewDosageReminderService creates a new DosageReminderService.
func NewDosageReminderService(
	storage DosageReminderStorage,
	notifs *notification.UserNotificationService,
	slog *slog.Logger,
	lc fx.Lifecycle,
) *DosageReminderService {
	s := &DosageReminderService{
		storage: storage,
		notifs:  notifs,
	}

	fakectx, stop := context.WithCancel(context.Background())
	done := make(chan struct{})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				s.run(fakectx, slog)
				close(done)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			stop()
			<-done
			return nil
		},
	})

	return nil
}

func (s *DosageReminderService) run(ctx context.Context, slog *slog.Logger) {
	nextRun := time.NewTimer(0)
	now := time.Now()

	for {
		slog.Debug("DosageReminderService: running update cycle")

		dosageRemindersIter := s.storage.UpcomingDosageReminders(ctx)

		tracked, err := ingestReminders(now, dosageRemindersIter, slog)
		if err != nil {
			slog.Error(
				"DosageReminderService: error ingesting reminders",
				"err", err)
			nextRun.Reset(nextUpdateIntervalOnError)
			goto skipToTimer
		}

		nextRun.Reset(tracked.nextRun.Sub(now))
		slog.Debug(
			"DosageReminderService: scheduling next run",
			"nextRun", tracked.nextRun)

		for _, r := range tracked.notifyingReminders {
			start := time.Now()
			err := s.notifs.NotifyUser(ctx, r.UserSecret, notificationapi.ReminderMessage)
			taken := time.Since(start)

			if err != nil {
				slog.ErrorContext(ctx,
					"DosageReminderService: error notifying user",
					"reminder.username", r.Username,
					"timeTaken", taken,
					"err", err)
			} else {
				slog.DebugContext(ctx,
					"DosageReminderService: notified user",
					"reminder.username", r.Username,
					"timeTaken", taken)
			}

			if err := s.storage.RecordRemindedDoses(ctx, []RemindedDose{r.toRemindedDose()}); err != nil {
				slog.ErrorContext(ctx,
					"DosageReminderService: error recording reminded doses",
					"reminder.username", r.Username,
					"err", err)
			}
		}

	skipToTimer:
		select {
		case <-ctx.Done():
			slog.Debug("DosageReminderService: stopping update cycle")
			return
		case now = <-nextRun.C:
			// keep running
		}
	}
}

type trackedDosageReminders struct {
	notifyingReminders []notifyingReminder
	nextRun            time.Time
}

type notifyingReminder struct {
	DosageReminder
	ClearSnooze bool
}

func (r notifyingReminder) toRemindedDose() RemindedDose {
	return RemindedDose{
		UserSecret:   r.UserSecret,
		RemindedDose: r.LastDose.TakenAt,
		ClearSnooze:  r.ClearSnooze,
	}
}

// ingestReminders ingests the streaming reminders into the tracker.
func ingestReminders(now time.Time, reminders iter.Seq2[DosageReminder, error], slog *slog.Logger) (*trackedDosageReminders, error) {
	notifyingReminders := make([]notifyingReminder, 0, 12)

	cutoffPoint := now.Add(nextUpdateInterval)
	earliestNextNotification := cutoffPoint

	for r, err := range reminders {
		if err != nil {
			return nil, err
		}

		nextNotification := r.NextNotification()
		if nextNotification.After(cutoffPoint) {
			slog.Debug(
				"ingestReminders: reminder is not relevant because it is too far into the future",
				"now", now,
				"reminder.nextNotification", nextNotification,
				"reminder.username", r.Username)
			continue
		}

		if r.LastRemindedDose != nil && r.SnoozedUntil == nil && r.LastRemindedDose.Equal(r.LastDose.TakenAt) {
			slog.Debug(
				"ingestReminders: reminder is not relevant because it was already reminded and not snoozed",
				"reminder.username", r.Username,
				"reminder.lastDose", r.LastDose.TakenAt)
			continue
		}

		if nextNotification.Before(now) {
			slog.Debug(
				"ingestReminders: recorded reminder for notification",
				"reminder.username", r.Username,
				"reminder.nextNotification", nextNotification)

			notifyingReminders = append(notifyingReminders, notifyingReminder{
				DosageReminder: r,
				ClearSnooze:    false, // TODO(diamondburned): implement snoozing
			})
			continue
		}

		if nextNotification.Before(earliestNextNotification) {
			slog.Debug(
				"ingestReminders: found new earliest next notification",
				"reminder.username", r.Username,
				"reminder.nextNotification", nextNotification)

			earliestNextNotification = nextNotification
			continue
		}
	}

	// Make sure this is not in the past.
	if earliestNextNotification.Before(now) {
		earliestNextNotification = now
	}

	// Make sure the next notification is at least 5 minutes in the future.
	if earliestNextNotification.Before(now.Add(shortestNextNotification)) {
		earliestNextNotification = now.Add(shortestNextNotification)
	}

	slog.Debug(
		"ingestReminders: figured out all relevant reminders",
		"ingestedAt", now,
		"timeTaken", time.Since(now),
		"numRelevantReminders", len(notifyingReminders),
		"earliestNextNotification", earliestNextNotification)

	return &trackedDosageReminders{
		notifyingReminders: notifyingReminders,
		nextRun:            earliestNextNotification,
	}, nil
}

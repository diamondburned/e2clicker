package dosage

import (
	"context"
	"log/slog"
	"time"

	"e2clicker.app/internal/ptr"
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
	UpcomingDosageReminders(ctx context.Context) ([]DosageReminder, error)

	// RecordRemindedDoseAttempts records the reminded dose attempts.
	// This is used to mark the reminder as sent or failed.
	RecordRemindedDoseAttempts(ctx context.Context, remindedDoses []RemindedDoseAttempt) error
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
	//
	// Deprecated: consider removing this field entirely, as it is not used in any
	// calculations anymore.
	LastRemindedDose *time.Time
	// LastRemindedAt is the time when the reminder was last sent, if it was at
	// all.
	LastRemindedAt *time.Time // sent_at
	// SnoozedUntil is the time until the reminder is snoozed.
	// This field is optional and is only set if the reminder is snoozed by the
	// user on the last notification.
	SnoozedUntil *time.Time
}

// RemindedDoseAttempt is a dose that has been reminded.
type RemindedDoseAttempt struct {
	// UserSecret is the secret of the user.
	UserSecret user.Secret
	// RemindedAt is the time when the reminder was sent.
	RemindedAt time.Time
	// RemindedDose is the dose that was reminded.
	// This is the TakenAt time of the dose.
	RemindedDose time.Time
	// ClearSnooze is true if the snooze should be cleared.
	// This is the case if the reminder was sent at the snoozed time.
	ClearSnooze bool
	// ErrorReason is the error reason if the reminder failed, if any.
	ErrorReason *string
}

// NextNotification returns the time of the next notification only if it should
// have been sent, meaning that the notification time is in the past of [now].
//
// It returns the snoozed time if the reminder is snoozed, otherwise it returns
// the time of the last dose plus the interval.
func (r DosageReminder) NextNotification() (time.Time, bool) {
	if r.SnoozedUntil != nil && !r.hasRemindedAtTime(*r.SnoozedUntil) {
		return *r.SnoozedUntil, true
	}

	if !r.LastDose.TakenAt.IsZero() {
		latestRemindedAt := r.LastDose.TakenAt.Add(r.Dosage.Interval.ToDuration())
		if !r.hasRemindedAtTime(latestRemindedAt) {
			return latestRemindedAt, true
		}

		for _, recurrence := range r.Dosage.ReminderRecurrence {
			recurrenceTime := latestRemindedAt.Add(recurrence.ToDuration())
			if !r.hasRemindedAtTime(recurrenceTime) {
				return recurrenceTime, true
			}
		}
	}

	return time.Time{}, false
}

func (r DosageReminder) hasRemindedAtTime(t time.Time) bool {
	return r.LastRemindedAt != nil && (r.LastRemindedAt.Equal(t) || r.LastRemindedAt.After(t))
}

const (
	nextUpdateInterval        = 30 * time.Minute
	nextUpdateIntervalOnError = 2 * time.Minute
)

// DosageReminderService is a service for managing dosage reminders.
type DosageReminderService struct {
	storage DosageReminderStorage
	notifs  notification.UserNotificationService
	logger  *slog.Logger
}

// NewDosageReminderService creates a new DosageReminderService.
func NewDosageReminderService(
	storage DosageReminderStorage,
	notifs notification.UserNotificationService,
	slog *slog.Logger,
	lc fx.Lifecycle,
) *DosageReminderService {
	s := &DosageReminderService{
		storage: storage,
		notifs:  notifs,
		logger:  slog,
	}

	fakectx, stop := context.WithCancel(context.Background())
	done := make(chan struct{})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				s.run(fakectx)
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

	return s
}

func (s *DosageReminderService) run(ctx context.Context) {
	nextRun := time.NewTimer(0)
	now := time.Now()

	for {
		nextRunAt := s.cycle(ctx, now)
		nextRun.Reset(time.Until(nextRunAt))

		slog.Debug(
			"DosageReminderService: scheduling next run",
			"nextRun", nextRun)

		select {
		case <-ctx.Done():
			slog.Debug("DosageReminderService: stopping update cycle")
			return
		case now = <-nextRun.C:
			// keep running
		}
	}
}

func (s *DosageReminderService) cycle(ctx context.Context, now time.Time) time.Time {
	slog := s.logger.With("now", now)
	slog.Debug("DosageReminderService: running update cycle")

	dosageReminders, err := s.storage.UpcomingDosageReminders(ctx)
	if err != nil {
		slog.Error(
			"DosageReminderService: error querying reminders",
			"err", err)
		return now.Add(nextUpdateIntervalOnError)
	}

	tracked, err := ingestReminders(now, dosageReminders, slog)
	if err != nil {
		slog.Error(
			"DosageReminderService: error ingesting reminders",
			"err", err)
		return now.Add(nextUpdateIntervalOnError)
	}

	for _, r := range tracked.notifyingReminders {
		start := time.Now()
		err := s.notifs.NotifyUser(ctx, r.UserSecret, notificationapi.ReminderMessage)
		taken := time.Since(start)

		attempt := RemindedDoseAttempt{
			UserSecret:   r.UserSecret,
			RemindedAt:   now,
			RemindedDose: r.LastDose.TakenAt,
			ClearSnooze:  r.ClearSnooze,
		}

		if err != nil {
			attempt.ErrorReason = ptr.To(err.Error())

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

		if err := s.storage.RecordRemindedDoseAttempts(ctx, []RemindedDoseAttempt{attempt}); err != nil {
			slog.ErrorContext(ctx,
				"DosageReminderService: error recording reminded doses",
				"reminder.username", r.Username,
				"err", err)
		}
	}

	return tracked.nextRun
}

type trackedDosageReminders struct {
	notifyingReminders []notifyingReminder
	nextRun            time.Time
}

type notifyingReminder struct {
	DosageReminder
	ClearSnooze bool
}

// ingestReminders ingests the streaming reminders into the tracker.
func ingestReminders(now time.Time, reminders []DosageReminder, slog *slog.Logger) (*trackedDosageReminders, error) {
	notifyingReminders := make([]notifyingReminder, 0, 12)

	for _, r := range reminders {
		nextNotification, ok := r.NextNotification()
		if !ok {
			slog.Debug(
				"ingestReminders: no upcoming notifications",
				"now", now,
				"reminder.nextNotification", nil,
				"reminder.username", r.Username)
			continue
		}

		if nextNotification.After(now) {
			slog.Debug(
				"ingestReminders: next notification is not due yet",
				"now", now,
				"reminder.username", r.Username,
				"reminder.nextNotification", nextNotification)
			continue
		}

		slog.Debug(
			"ingestReminders: recorded reminder for next notification",
			"now", now,
			"reminder.username", r.Username,
			"reminder.nextNotification", nextNotification)

		notifyingReminders = append(notifyingReminders, notifyingReminder{
			DosageReminder: r,
			ClearSnooze:    false, // TODO(diamondburned): implement snoozing
		})
	}

	slog.Debug(
		"ingestReminders: figured out all relevant reminders",
		"now", now,
		"timeTaken", time.Since(now),
		"numRelevantReminders", len(notifyingReminders))

	return &trackedDosageReminders{
		notifyingReminders: notifyingReminders,
		nextRun:            now.Add(nextUpdateInterval),
	}, nil
}

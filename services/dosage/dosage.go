package dosage

import (
	"context"
	"errors"
	"iter"
	"time"

	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/e2clicker/services/dosage/openapi"
	"libdb.so/e2clicker/services/user"
)

func init() {
	publicerrors.MarkValuesPublic(
		ErrNoDoseMatched,
	)
}

var (
	ErrNoDoseMatched = errors.New("no dose matched")
)

// DosageStorage is a storage for dosage data.
type DosageStorage interface {
	// DeliveryMethods returns the available delivery methods.
	DeliveryMethods(ctx context.Context) ([]DeliveryMethod, error)
	// Dosage returns the dosage for a user.
	// If the user has no schedule yet, this returns nil.
	Dosage(ctx context.Context, secret user.Secret) (*Dosage, error)
	// SetDosage sets the dosage for a user.
	// The user secret is taken from the Schedule.
	SetDosage(ctx context.Context, s Dosage) error
	// ClearDosage clears the dosage for a user.
	ClearDosage(ctx context.Context, secret user.Secret) error
}

// DoseHistoryStorage is a storage for dose history data.
type DoseHistoryStorage interface {
	// RecordDose records a single dose. The dose observation's ID is returned.
	RecordDose(ctx context.Context, secret user.Secret, dose Dose) (int64, error)
	// ImportDoses imports doses in bulk and returns the number of doses
	// imported.
	// This is a separate method from RecordDose to allow for more efficient
	// bulk imports.
	ImportDoses(ctx context.Context, secret user.Secret, doseSeq iter.Seq2[Dose, error]) (int64, error)
	// EditDose edits a dose by ID. All fields are updated.
	EditDose(ctx context.Context, secret user.Secret, doseIDs int64, dose Dose) error
	// ForgetDoses forgets the given doses.
	ForgetDoses(ctx context.Context, secret user.Secret, doseIDs []int64) error
	// DoseHistory returns the history of a dosage schedule. If end is zero, it
	// is considered to be now. If begin is zero, it is considered to be
	// since the beginning of time.
	// The history is ordered by time taken, with the oldest dose first.
	// If there's an error, the returned sequence will yield the error with a
	// zero-value [Observation].
	DoseHistory(ctx context.Context, secret user.Secret, begin, end time.Time) iter.Seq2[Observation, error]
}

// RecordedDosesResult is the result of recording doses.
type RecordedDosesResult struct {
	// Created is the number of doses that were created.
	// This is the number of doses that were not already in the database.
	Created int
}

// DeliveryMethod describes a method of delivery for medication.
type DeliveryMethod = openapi.DeliveryMethod

// Dosage describes a dosage schedule.
type Dosage struct {
	// UserSecret is the secret of the user who the schedule is for.
	UserSecret user.Secret
	// DeliveryMethod is the method of delivery for the medication.
	// Check the [delivery_methods] table.
	DeliveryMethod string
	// Dose is the amount of medication to be delivered/taken.
	Dose float32
	// Interval is the interval between doses in days.
	Interval Days
	// Concurrence is the number of estrogen patches that are on the body at
	// once. This is only relevant if DeliveryMethod is "patch".
	Concurrence *int
}

// Days is a number of days. It acts as a duration of time, so 1.5 Days is
// 36 hours.
type Days float64

// ToDuration converts Days to time.Duration.
func (d Days) ToDuration() time.Duration {
	return time.Duration(float64(d) * float64(24*time.Hour))
}

// Dose describes a dose of medication in time.
type Dose struct {
	// DeliveryMethod is the method of delivery for the medication
	// at the time the dose was taken.
	DeliveryMethod string
	// Dose is the amount of medication that was taken.
	Dose float32
	// TakenAt is the time the dose was taken.
	TakenAt time.Time
	// TakenOffAt is the time the dose was taken off the body.
	// This is only relevant if DeliveryMethod is "patch".
	TakenOffAt *time.Time
	// Comment is a comment about the dose.
	Comment string
}

// Observation is an observed dose of medication. It is used in the history of a
// dosage schedule.
type Observation struct {
	// ID is the ID of the observation.
	ID int64
	Dose
}

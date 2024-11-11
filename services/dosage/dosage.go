package dosage

import (
	"context"
	"errors"
	"time"

	"libdb.so/e2clicker/internal/publicerrors"
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
	// RecordDose records a dose.
	RecordDose(ctx context.Context, secret user.Secret, takenAt time.Time) (Observation, error)
	// EditDose edits a dose.
	EditDose(ctx context.Context, secret user.Secret, o Observation) error
	// ForgetDoses forgets the given doses.
	ForgetDoses(ctx context.Context, secret user.Secret, doseIDs []int64) error
	// DoseHistory returns the history of a dosage schedule.
	DoseHistory(ctx context.Context, secret user.Secret, begin, end time.Time) (History, error)
}

// DeliveryMethod describes a method of delivery for medication.
type DeliveryMethod struct {
	// ID is a short string representing the delivery method.
	// This is what goes into the DeliveryMethod fields.
	ID string
	// Units is the units of the dose.
	Units string
	// Name is the full name of the delivery method.
	Name string
}

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

// History is the history of a dosage schedule.
type History struct {
	// UserSecret is the secret of the user who the schedule is for.
	UserSecret user.Secret
	// Entries is the dosage data over time.
	Entries []Observation
}

// Observation is a point in the history of a dosage schedule.
type Observation struct {
	// DoseID is the ID of the dose.
	DoseID int64
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
}

package dosage

import (
	"context"
	"errors"
	"fmt"
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
	RecordDose(ctx context.Context, secret user.Secret, dose Dose) error
	// ImportDoses imports doses in bulk and returns the number of doses
	// imported.
	// This is a separate method from RecordDose to allow for more efficient
	// bulk imports.
	ImportDoses(ctx context.Context, secret user.Secret, doseSeq iter.Seq[Dose]) (int64, error)
	// EditDose edits a dose by the previous takenAt time.
	// All fields are updated.
	EditDose(ctx context.Context, secret user.Secret, doseTime time.Time, dose Dose) error
	// ForgetDoses forgets the given doses.
	ForgetDoses(ctx context.Context, secret user.Secret, doseTimes []time.Time) error
	// DoseHistory returns the history of a dosage schedule. If end is zero, it
	// is considered to be now. If begin is zero, it is considered to be
	// since the beginning of time.
	// The history is ordered by time taken, with the oldest dose first.
	// If there's an error, the returned sequence will yield the error with a
	// zero-value [Observation].
	DoseHistory(ctx context.Context, secret user.Secret, begin, end time.Time) iter.Seq2[Dose, error]
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

func doseFromOpenAPI(d openapi.Dose) Dose {
	var comment string
	if d.Comment != nil {
		comment = *d.Comment
	}

	return Dose{
		DeliveryMethod: d.DeliveryMethod,
		Dose:           d.Dose,
		TakenAt:        d.TakenAt,
		TakenOffAt:     d.TakenOffAt,
		Comment:        comment,
	}
}

func (d Dose) ToOpenAPI() openapi.Dose {
	var commentPtr *string
	if d.Comment != "" {
		commentPtr = &d.Comment
	}

	return openapi.Dose{
		DeliveryMethod: d.DeliveryMethod,
		Dose:           d.Dose,
		TakenAt:        d.TakenAt,
		TakenOffAt:     d.TakenOffAt,
		Comment:        commentPtr,
	}
}

// CSVDoseRecord is a single dose record in a CSV file.
// This is version 1.
type CSVDoseRecord struct {
	TakenAt        string  `csv:"takenAt"`    // RFC3339
	TakenOffAt     string  `csv:"takenOffAt"` // RFC3339 or ""
	DeliveryMethod string  `csv:"deliveryMethod"`
	Dose           float32 `csv:"dose"`
	Comment        string  `csv:"comment"`
}

func doseFromCSV(r CSVDoseRecord) (Dose, error) {
	takenAt, err := parseCSVTime(r.TakenAt)
	if err != nil {
		return Dose{}, fmt.Errorf("cannot parse takenAt: %w", err)
	}

	var takenOffAt *time.Time
	if r.TakenOffAt != "" {
		t, err := parseCSVTime(r.TakenOffAt)
		if err != nil {
			return Dose{}, fmt.Errorf("cannot parse takenOffAt: %w", err)
		}
		takenOffAt = &t
	}

	return Dose{
		Dose:           r.Dose,
		DeliveryMethod: r.DeliveryMethod,
		TakenAt:        takenAt,
		TakenOffAt:     takenOffAt,
		Comment:        r.Comment,
	}, nil
}

func (d Dose) ToCSV() CSVDoseRecord {
	return CSVDoseRecord{
		Dose:           d.Dose,
		DeliveryMethod: d.DeliveryMethod,
		TakenAt:        formatCSVTime(d.TakenAt),
		TakenOffAt:     optionalStr(d.TakenOffAt, formatCSVTime),
		Comment:        d.Comment,
	}
}

func optionalStr[T any](v *T, f func(T) string) string {
	if v == nil {
		return ""
	}
	return f(*v)
}

func formatCSVTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func parseCSVTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse time %q: expected RFC3339, got %w", s, err)
	}
	return t, nil
}

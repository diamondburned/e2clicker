package dosage

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"time"
	"uuid"

	"e2clicker.app/internal/ptr"
	"e2clicker.app/internal/publicerrors"
	"e2clicker.app/services/dosage/openapi"
	"e2clicker.app/services/user"
	openapi_types "github.com/oapi-codegen/runtime/types"
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
	SetDosage(ctx context.Context, secret user.Secret, s Dosage) error
	// ClearDosage clears the dosage for a user.
	ClearDosage(ctx context.Context, secret user.Secret) error
}

// DoseHistoryStorage is a storage for dose history data.
type DoseHistoryStorage interface {
	// RecordDose records a single dose and returns the recorded observation
	// including its assigned ID.
	RecordDose(ctx context.Context, secret user.Secret, dose Dose) (Dose, error)
	// ImportDoses imports doses in bulk and returns the number of doses
	// imported.
	// This is a separate method from RecordDose to allow for more efficient
	// bulk imports.
	ImportDoses(ctx context.Context, secret user.Secret, doseSeq iter.Seq[Dose]) (int64, error)
	// EditDose edits a dose by its ID.
	// All fields are updated.
	EditDose(ctx context.Context, secret user.Secret, id uuid.UUID, dose Dose) error
	// ForgetDoses forgets the given doses by their IDs.
	ForgetDoses(ctx context.Context, secret user.Secret, ids []uuid.UUID) error
	// DoseHistory returns the history of a dosage schedule. If end is zero, it
	// is considered to be now. If begin is zero, it is considered to be
	// since the beginning of time.
	// The history is ordered by time taken, with the oldest dose first.
	// If there's an error, the returned sequence will yield the error with a
	// zero-value [Observation].
	DoseHistory(ctx context.Context, secret user.Secret, begin, end time.Time) ([]Dose, error)
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
type Dosage = openapi.Dosage

// Dose describes a dose of medication in time.
type Dose struct {
	// ID is the unique identifier (UUIDv7) of the dose observation.
	ID uuid.UUID
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
	return Dose{
		ID:             uuid.UUID(d.ID),
		DeliveryMethod: d.DeliveryMethod,
		Dose:           d.Dose,
		TakenAt:        d.TakenAt,
		TakenOffAt:     d.TakenOffAt,
		Comment:        ptr.Deref(d.Comment),
	}
}

func doseFromImportOpenAPI(d openapi.DoseImport) Dose {
	var id uuid.UUID
	if d.ID != nil {
		id = uuid.UUID(*d.ID)
	}
	return Dose{
		ID:             id,
		DeliveryMethod: d.DeliveryMethod,
		Dose:           d.Dose,
		TakenAt:        d.TakenAt,
		TakenOffAt:     d.TakenOffAt,
		Comment:        ptr.Deref(d.Comment),
	}
}

func (d Dose) ToOpenAPI() openapi.Dose {
	var commentPtr *string
	if d.Comment != "" {
		commentPtr = &d.Comment
	}

	return openapi.Dose{
		ID:             openapi_types.UUID(d.ID),
		DeliveryMethod: d.DeliveryMethod,
		Dose:           d.Dose,
		TakenAt:        d.TakenAt,
		TakenOffAt:     d.TakenOffAt,
		Comment:        commentPtr,
	}
}

// CSVDoseRecord is a single dose record in an e2clicker v2 CSV file.
type CSVDoseRecord struct {
	ID             string  `csv:"id"` // UUID or ""
	DeliveryMethod string  `csv:"deliveryMethod"`
	Dose           float32 `csv:"dose"`
	TakenAt        string  `csv:"takenAt"`    // RFC3339
	TakenOffAt     string  `csv:"takenOffAt"` // RFC3339 or ""
	Comment        string  `csv:"comment"`
}

// csvDoseRecordV1 is a single dose record in a legacy e2clicker v1 CSV file.
type csvDoseRecordV1 struct {
	DeliveryMethod string  `csv:"deliveryMethod"`
	Dose           float32 `csv:"dose"`
	TakenAt        string  `csv:"takenAt"`    // RFC3339
	TakenOffAt     string  `csv:"takenOffAt"` // RFC3339 or ""
	Comment        string  `csv:"comment"`
}

// EstrannaiseCSVDoseRecord is a single dose record in estrannaise.js CSV format.
// This format is export-only and cannot be imported.
type EstrannaiseCSVDoseRecord struct {
	Dose  float32 `csv:"dose"`
	Days  float64 `csv:"days"`
	Model string  `csv:"model"`
}

func doseFromCSV(r CSVDoseRecord) (Dose, error) {
	var id uuid.UUID
	if r.ID != "" {
		var err error
		id, err = uuid.Parse(r.ID)
		if err != nil {
			return Dose{}, publicerrors.Errorf("cannot parse id %q: %w", r.ID, err)
		}
	}

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
		ID:             id,
		DeliveryMethod: r.DeliveryMethod,
		Dose:           r.Dose,
		TakenAt:        takenAt,
		TakenOffAt:     takenOffAt,
		Comment:        r.Comment,
	}, nil
}

func doseFromCSVV1(r csvDoseRecordV1) (Dose, error) {
	return doseFromCSV(CSVDoseRecord{
		DeliveryMethod: r.DeliveryMethod,
		Dose:           r.Dose,
		TakenAt:        r.TakenAt,
		TakenOffAt:     r.TakenOffAt,
		Comment:        r.Comment,
	})
}

func (d Dose) ToCSV() CSVDoseRecord {
	var id string
	if d.ID != uuid.Nil() {
		id = d.ID.String()
	}

	return CSVDoseRecord{
		ID:             id,
		DeliveryMethod: d.DeliveryMethod,
		Dose:           d.Dose,
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

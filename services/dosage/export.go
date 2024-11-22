package dosage

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"go.uber.org/fx"
	"golang.org/x/time/rate"
	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/e2clicker/internal/userlimit"
	"libdb.so/e2clicker/services/user"
	"libdb.so/xcsv"
)

// Allow waiting for 500ms max before returning [ExportLimitError].
const maxDelay = 500 * time.Millisecond

// ExportFormat is the format of the export.
type ExportFormat string

const (
	ExportCSV  ExportFormat = "text/csv"
	ExportJSON ExportFormat = "application/json"
)

// ExporterService exports dosage data to various formats.
type ExporterService struct {
	storage DoseHistoryStorage
	logger  *slog.Logger

	importLimiter *userlimit.UserRateLimiter[user.Secret]
	exportLimiter *userlimit.UserRateLimiter[user.Secret]
}

// NewExporterService creates a new CSVExporterService.
func NewExporterService(storage DoseHistoryStorage, lc fx.Lifecycle, logger *slog.Logger) *ExporterService {
	s := &ExporterService{
		storage:       storage,
		logger:        logger,
		importLimiter: userlimit.NewUserRateLimiter[user.Secret](rate.Every(15*time.Minute), 3),
		exportLimiter: userlimit.NewUserRateLimiter[user.Secret](rate.Every(15*time.Minute), 3),
	}

	stopCleanup := s.importLimiter.BeginCleanup()
	lc.Append(fx.StopHook(func(ctx context.Context) error {
		stopCleanup()
		return nil
	}))

	return s
}

/*
 * CSV Export
 */

// CSVDoseRecord is a single dose record in a CSV file.
// This is version 1.
type CSVDoseRecord struct {
	TakenAt        string  `csv:"takenAt"`    // RFC3339
	TakenOffAt     string  `csv:"takenOffAt"` // RFC3339 or ""
	DeliveryMethod string  `csv:"deliveryMethod"`
	Dose           float32 `csv:"dose"`
	Comment        string  `csv:"comment"`
}

var csvDoseRecordColumns = xcsv.ColumnNames[CSVDoseRecord]()

// ExportDoseHistoryOptions are options for exporting dose history as a CSV
// file.
type ExportDoseHistoryOptions struct {
	Begin  time.Time // zero means from beginning of time
	End    time.Time // zero means until now
	Format ExportFormat
}

// ExportDoseHistory exports the dose history of the user as a CSV file into the
// given writer. It returns the number of records exported.
func (s *ExporterService) ExportDoseHistory(ctx context.Context, out io.Writer, secret user.Secret, o ExportDoseHistoryOptions) (int64, error) {
	if err := userlimit.AsError(s.exportLimiter.Reserve(secret)); err != nil {
		return 0, err
	}

	var exported int64
	var scanErrs []error
	doseIter := func(yield func(Dose) bool) {
		for o, err := range s.storage.DoseHistory(ctx, secret, o.Begin, o.End) {
			if err != nil {
				scanErrs = append(scanErrs, err)
				continue
			}
			if !yield(o.Dose) {
				break
			}
			exported++
		}
	}

	var err error

	switch o.Format {
	case ExportCSV:
		csvw := csv.NewWriter(out)
		csvw.Write(csvDoseRecordColumns)

		err = xcsv.Marshal(csvw, func(yield func(CSVDoseRecord) bool) {
			for dose := range doseIter {
				if !yield(CSVDoseRecord{
					Dose:           dose.Dose,
					DeliveryMethod: dose.DeliveryMethod,
					TakenAt:        formatCSVTime(dose.TakenAt),
					TakenOffAt:     optionalStr(dose.TakenOffAt, formatCSVTime),
					Comment:        dose.Comment,
				}) {
					break
				}
			}
		})
	case ExportJSON:
		return 0, publicerrors.New("JSON export is not supported yet")
	default:
		return 0, publicerrors.Errorf("unsupported export format %q", o.Format)
	}

	if err != nil {
		return exported, err
	}

	return exported, errors.Join(scanErrs...)
}

// ImportDoseHistoryOptions are options for importing dose history from a file.
type ImportDoseHistoryOptions struct {
	Format ExportFormat
}

// ImportDoseHistoryResult is the result of importing dose history from a file.
type ImportDoseHistoryResult struct {
	Records   int64
	Succeeded int64
}

// ImportDoseHistory imports dose history from a CSV file.
func (s *ExporterService) ImportDoseHistory(ctx context.Context, in io.Reader, secret user.Secret, o ImportDoseHistoryOptions) (ImportDoseHistoryResult, error) {
	if err := userlimit.AsError(s.importLimiter.Reserve(secret)); err != nil {
		return ImportDoseHistoryResult{}, err
	}

	var doseIter func(yield func(Dose, error) bool)

	switch o.Format {
	case ExportCSV:
		csvr := csv.NewReader(in)
		csvr.Comment = '#'
		csvr.LazyQuotes = true
		csvr.ReuseRecord = true
		csvr.TrimLeadingSpace = true

		doseIter = func(yield func(Dose, error) bool) {
			for r, err := range xcsv.Unmarshal[CSVDoseRecord](csvr, xcsv.SkipHeader()) {
				if err != nil {
					if !yield(Dose{}, err) {
						return
					}
					continue
				}

				takenAt, err := parseCSVTime(r.TakenAt)
				if err != nil {
					if !yield(Dose{}, err) {
						return
					}
					continue
				}

				var takenOffAt *time.Time
				if r.TakenOffAt != "" {
					t, err := parseCSVTime(r.TakenOffAt)
					if err != nil {
						if !yield(Dose{}, err) {
							return
						}
						continue
					}
					takenOffAt = &t
				}

				dose := Dose{
					Dose:           r.Dose,
					DeliveryMethod: r.DeliveryMethod,
					TakenAt:        takenAt,
					TakenOffAt:     takenOffAt,
					Comment:        r.Comment,
				}
				if !yield(dose, nil) {
					return
				}
			}
		}
	case ExportJSON:
		return ImportDoseHistoryResult{}, publicerrors.New("JSON import is not supported yet")
	default:
		return ImportDoseHistoryResult{}, publicerrors.Errorf("unsupported import format %q", o.Format)
	}

	var records int64
	var importErrors []error

	// ImportDoses does allow us to yield an error to stop the whole insertion,
	// but we want to import as many records as possible, so we collect errors
	// and return them at the end.
	succeeded, err := s.storage.ImportDoses(ctx, secret, func(yield func(Dose, error) bool) {
		for dose, err := range doseIter {
			if err != nil {
				importErrors = append(importErrors, err)
				continue
			}

			records++
			if !yield(dose, nil) {
				return
			}
		}
	})

	r := ImportDoseHistoryResult{
		Records:   records,
		Succeeded: succeeded,
	}

	if err != nil {
		return r, err
	}

	return r, errors.Join(importErrors...)
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

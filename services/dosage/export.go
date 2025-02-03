package dosage

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"iter"
	"log/slog"
	"time"

	"go.uber.org/fx"
	"golang.org/x/time/rate"
	"e2clicker.app/internal/jsonarray"
	"e2clicker.app/internal/publicerrors"
	"e2clicker.app/internal/userlimit"
	"e2clicker.app/services/dosage/openapi"
	"e2clicker.app/services/user"
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

// AsMIME returns itself.
func (f ExportFormat) AsMIME() string { return string(f) }

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
	limit := s.exportLimiter.Reserve(secret)
	if err := userlimit.AsError(limit); err != nil {
		return 0, err
	}

	var exported int64
	var scanErrs []error
	history := func(yield func(Dose) bool) {
		for o, err := range s.storage.DoseHistory(ctx, secret, o.Begin, o.End) {
			if err != nil {
				scanErrs = append(scanErrs, err)
				continue
			}
			if !yield(o) {
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
			for o := range history {
				if !yield(o.ToCSV()) {
					break
				}
			}
		})
	case ExportJSON:
		err = jsonarray.MarshalArray(out, func(yield func(openapi.Dose) bool) {
			for o := range history {
				if !yield(o.ToOpenAPI()) {
					break
				}
			}
		})
	default:
		err = publicerrors.Errorf("unsupported export format %q", o.Format)
	}

	if err != nil {
		limit.Cancel()
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
	limit := s.importLimiter.Reserve(secret)
	if err := userlimit.AsError(limit); err != nil {
		return ImportDoseHistoryResult{}, err
	}

	var doses iter.Seq[Dose]
	var records int64
	var importErrors []error

	switch o.Format {
	case ExportCSV:
		csvr := csv.NewReader(in)
		csvr.Comment = '#'
		csvr.LazyQuotes = true
		csvr.ReuseRecord = true
		csvr.TrimLeadingSpace = true

		doses = func(yield func(Dose) bool) {
			for r, err := range xcsv.Unmarshal[CSVDoseRecord](csvr, xcsv.SkipHeader()) {
				if err != nil {
					importErrors = append(importErrors, err)
					continue
				}

				d, err := doseFromCSV(r)
				if err != nil {
					importErrors = append(importErrors, err)
					continue
				}

				if !yield(d) {
					return
				}

				records++
			}
		}
	case ExportJSON:
		doses = func(yield func(Dose) bool) {
			for r, err := range jsonarray.UnmarshalArray[openapi.Dose](in) {
				if err != nil {
					importErrors = append(importErrors, err)
					continue
				}

				if !yield(doseFromOpenAPI(r)) {
					return
				}

				records++
			}
		}
	default:
		limit.Cancel()
		return ImportDoseHistoryResult{}, publicerrors.Errorf("unsupported import format %q", o.Format)
	}

	succeeded, err := s.storage.ImportDoses(ctx, secret, doses)
	r := ImportDoseHistoryResult{
		Records:   records,
		Succeeded: succeeded,
	}

	if err != nil {
		limit.Cancel()
		return r, err
	}

	return r, errors.Join(importErrors...)
}

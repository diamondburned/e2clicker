package dosage

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"iter"
	"log/slog"
	"slices"
	"strconv"
	"time"

	"e2clicker.app/internal/jsonarray"
	"e2clicker.app/internal/publicerrors"
	"e2clicker.app/internal/userlimit"
	"e2clicker.app/services/dosage/openapi"
	"e2clicker.app/services/user"
	"go.uber.org/fx"
	"golang.org/x/time/rate"
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

// CSVFormat selects the CSV schema when exporting as text/csv.
type CSVFormat string

const (
	CSVFormatE2Clicker   CSVFormat = "e2clicker"
	CSVFormatEstrannaise CSVFormat = "estrannaise"
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

var (
	csvDoseRecordColumns   = xcsv.ColumnNames[CSVDoseRecord]()
	csvDoseRecordV1Columns = xcsv.ColumnNames[csvDoseRecordV1]()
	estrannaiseCSVColumns  = xcsv.ColumnNames[EstrannaiseCSVDoseRecord]()
)

// ExportDoseHistoryOptions are options for exporting dose history as a CSV
// file.
type ExportDoseHistoryOptions struct {
	Begin     time.Time // zero means from beginning of time
	End       time.Time // zero means until now
	Format    ExportFormat
	CSVFormat CSVFormat // only applies to ExportCSV; defaults to CSVFormatE2Clicker
}

// ExportDoseHistory exports the dose history of the user as a CSV file into the
// given writer. It returns the number of records exported.
func (s *ExporterService) ExportDoseHistory(ctx context.Context, out io.Writer, secret user.Secret, o ExportDoseHistoryOptions) (int64, error) {
	if o.Format == ExportJSON && o.CSVFormat == CSVFormatEstrannaise {
		return 0, publicerrors.New("estrannaise format is only supported for text/csv exports")
	}

	limit := s.exportLimiter.Reserve(secret)
	if err := userlimit.AsError(limit); err != nil {
		return 0, err
	}

	doses, err := s.storage.DoseHistory(ctx, secret, o.Begin, o.End)
	if err != nil {
		return 0, err
	}

	var exported int64
	history := func(yield func(Dose) bool) {
		for _, o := range doses {
			if !yield(o) {
				break
			}
			exported++
		}
	}

	switch o.Format {
	case ExportCSV:
		csvw := csv.NewWriter(out)
		switch o.CSVFormat {
		case "", CSVFormatE2Clicker:
			csvw.Write(csvDoseRecordColumns)

			err = xcsv.Marshal(csvw, func(yield func(CSVDoseRecord) bool) {
				for o := range history {
					if !yield(o.ToCSV()) {
						break
					}
				}
			})
		case CSVFormatEstrannaise:
			csvw.Write(estrannaiseCSVColumns)

			var firstTakenAt time.Time
			err = xcsv.Marshal(csvw, func(yield func(EstrannaiseCSVDoseRecord) bool) {
				for o := range history {
					if firstTakenAt.IsZero() {
						firstTakenAt = o.TakenAt
					}
					days := o.TakenAt.Sub(firstTakenAt).Seconds() / 86400
					if !yield(EstrannaiseCSVDoseRecord{
						Dose:  o.Dose,
						Days:  days,
						Model: o.DeliveryMethod,
					}) {
						break
					}
				}
			})
		default:
			err = publicerrors.Errorf("unsupported CSV format %q", o.CSVFormat)
		}
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

	return exported, nil
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
			for d := range readCSVDoses(csvr, &importErrors) {
				if !yield(d) {
					return
				}

				records++
			}
		}
	case ExportJSON:
		doses = func(yield func(Dose) bool) {
			for r, err := range jsonarray.UnmarshalArray[openapi.DoseImport](in) {
				if err != nil {
					importErrors = append(importErrors, err)
					continue
				}

				if !yield(doseFromImportOpenAPI(r)) {
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

func readCSVDoses(csvr *csv.Reader, importErrors *[]error) iter.Seq[Dose] {
	return func(yield func(Dose) bool) {
		const (
			formatUnknown = iota
			formatV1
			formatV2
		)

		format := formatUnknown
		for line := 0; ; line++ {
			record, err := csvr.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				*importErrors = append(*importErrors, publicerrors.Errorf("failed to read CSV record at line %d: %w", line+1, err))
				continue
			}

			if line == 0 {
				switch {
				case slices.Equal(record, estrannaiseCSVColumns):
					*importErrors = append(*importErrors, publicerrors.New("estrannaise CSV format cannot be imported; use e2clicker CSV or JSON format instead"))
					return
				case slices.Equal(record, csvDoseRecordColumns):
					format = formatV2
					continue
				case slices.Equal(record, csvDoseRecordV1Columns):
					format = formatV1
					continue
				case len(record) == len(csvDoseRecordV1Columns):
					format = formatV1
				case len(record) == len(csvDoseRecordColumns):
					format = formatV2
				default:
					*importErrors = append(*importErrors, publicerrors.Errorf("line %d: expected %d (v1) or %d (v2) columns, got %d", line+1, len(csvDoseRecordV1Columns), len(csvDoseRecordColumns), len(record)))
					return
				}
			}

			var dose Dose
			switch format {
			case formatV1:
				if len(record) < len(csvDoseRecordV1Columns) {
					*importErrors = append(*importErrors, publicerrors.Errorf("line %d: wanted %d columns, got %d", line+1, len(csvDoseRecordV1Columns), len(record)))
					continue
				}
				doseAmt, err := strconv.ParseFloat(record[1], 32)
				if err != nil {
					*importErrors = append(*importErrors, publicerrors.Errorf("line %d: invalid dose %q: %w", line+1, record[1], err))
					continue
				}
				dose, err = doseFromCSVV1(csvDoseRecordV1{
					DeliveryMethod: record[0],
					Dose:           float32(doseAmt),
					TakenAt:        record[2],
					TakenOffAt:     record[3],
					Comment:        record[4],
				})
				if err != nil {
					*importErrors = append(*importErrors, publicerrors.Errorf("line %d: %w", line+1, err))
					continue
				}
			case formatV2:
				if len(record) < len(csvDoseRecordColumns) {
					*importErrors = append(*importErrors, publicerrors.Errorf("line %d: wanted %d columns, got %d", line+1, len(csvDoseRecordColumns), len(record)))
					continue
				}
				doseAmt, err := strconv.ParseFloat(record[2], 32)
				if err != nil {
					*importErrors = append(*importErrors, publicerrors.Errorf("line %d: invalid dose %q: %w", line+1, record[2], err))
					continue
				}
				dose, err = doseFromCSV(CSVDoseRecord{
					ID:             record[0],
					DeliveryMethod: record[1],
					Dose:           float32(doseAmt),
					TakenAt:        record[3],
					TakenOffAt:     record[4],
					Comment:        record[5],
				})
				if err != nil {
					*importErrors = append(*importErrors, publicerrors.Errorf("line %d: %w", line+1, err))
					continue
				}
			}

			if !yield(dose) {
				return
			}
		}
	}
}

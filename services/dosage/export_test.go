package dosage

import (
	"bytes"
	"context"
	"iter"
	"slices"
	"strings"
	"testing"
	"time"
	"uuid"

	"e2clicker.app/services/user"
	"github.com/alecthomas/assert/v2"
	"github.com/neilotoole/slogt"
	"go.uber.org/fx/fxtest"
)

type stubDoseHistoryStorage struct {
	history  []Dose
	imported []Dose
}

func (s *stubDoseHistoryStorage) RecordDose(_ context.Context, _ user.Secret, dose Dose) (Dose, error) {
	if dose.ID == uuid.Nil() {
		dose.ID = uuid.NewV7()
	}
	s.history = append(s.history, dose)
	return dose, nil
}

func (s *stubDoseHistoryStorage) ImportDoses(_ context.Context, _ user.Secret, doseSeq iter.Seq[Dose]) (int64, error) {
	for d := range doseSeq {
		s.imported = append(s.imported, d)
	}
	return int64(len(s.imported)), nil
}

func (s *stubDoseHistoryStorage) EditDose(_ context.Context, _ user.Secret, _ uuid.UUID, _ Dose) error {
	return nil
}

func (s *stubDoseHistoryStorage) ForgetDoses(_ context.Context, _ user.Secret, _ []uuid.UUID) error {
	return nil
}

func (s *stubDoseHistoryStorage) DoseHistory(_ context.Context, _ user.Secret, _, _ time.Time) ([]Dose, error) {
	return slices.Clone(s.history), nil
}

func TestExporterServiceCSVAndJSON(t *testing.T) {
	ctx := context.Background()
	id1 := uuid.MustParse("019450a0-0000-7000-8000-000000000001")
	id2 := uuid.MustParse("019450a0-0000-7000-8000-000000000002")
	t0 := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	t1 := t0.Add(84 * time.Hour) // 3.5 days

	t.Run("export_e2clicker_csv", func(t *testing.T) {
		storage := &stubDoseHistoryStorage{
			history: []Dose{
				{
					ID:             id1,
					DeliveryMethod: "EV im",
					Dose:           5,
					TakenAt:        t0,
					Comment:        "first",
				},
				{
					ID:             id2,
					DeliveryMethod: "patch",
					Dose:           0.1,
					TakenAt:        t1,
					TakenOffAt:     new(t1.Add(72 * time.Hour)),
				},
			},
		}
		svc := NewExporterService(storage, fxtest.NewLifecycle(t), slogt.New(t))

		var buf bytes.Buffer
		n, err := svc.ExportDoseHistory(ctx, &buf, "user-1", ExportDoseHistoryOptions{
			Format: ExportCSV,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), n)
		assert.Equal(t,
			"id,deliveryMethod,dose,takenAt,takenOffAt,comment\n"+
				"019450a0-0000-7000-8000-000000000001,EV im,5,2025-01-10T12:00:00Z,,first\n"+
				"019450a0-0000-7000-8000-000000000002,patch,0.1,2025-01-14T00:00:00Z,2025-01-17T00:00:00Z,\n",
			buf.String(),
		)
	})

	t.Run("export_estrannaise_csv", func(t *testing.T) {
		storage := &stubDoseHistoryStorage{
			history: []Dose{
				{ID: id1, DeliveryMethod: "EV im", Dose: 5, TakenAt: t0},
				{ID: id2, DeliveryMethod: "EV im", Dose: 6, TakenAt: t1},
			},
		}
		svc := NewExporterService(storage, fxtest.NewLifecycle(t), slogt.New(t))

		var buf bytes.Buffer
		n, err := svc.ExportDoseHistory(ctx, &buf, "user-1", ExportDoseHistoryOptions{
			Format:    ExportCSV,
			CSVFormat: CSVFormatEstrannaise,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), n)
		assert.Equal(t,
			"dose,days,model\n"+
				"5,0,EV im\n"+
				"6,3.5,EV im\n",
			buf.String(),
		)
	})

	t.Run("reject_estrannaise_json_export", func(t *testing.T) {
		svc := NewExporterService(&stubDoseHistoryStorage{}, fxtest.NewLifecycle(t), slogt.New(t))

		var buf bytes.Buffer
		_, err := svc.ExportDoseHistory(ctx, &buf, "user-1", ExportDoseHistoryOptions{
			Format:    ExportJSON,
			CSVFormat: CSVFormatEstrannaise,
		})
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "estrannaise"))
	})

	t.Run("import_v2_csv", func(t *testing.T) {
		storage := &stubDoseHistoryStorage{}
		svc := NewExporterService(storage, fxtest.NewLifecycle(t), slogt.New(t))

		res, err := svc.ImportDoseHistory(ctx, strings.NewReader(
			"id,deliveryMethod,dose,takenAt,takenOffAt,comment\n"+
				"019450a0-0000-7000-8000-000000000001,EV im,5,2025-01-10T12:00:00Z,,with id\n"+
				",EV im,4,2025-01-14T00:00:00Z,,without id\n",
		), "user-1", ImportDoseHistoryOptions{Format: ExportCSV})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), res.Records)
		assert.Equal(t, int64(2), res.Succeeded)
		assert.Equal(t, id1, storage.imported[0].ID)
		assert.Equal(t, "with id", storage.imported[0].Comment)
		assert.Equal(t, uuid.Nil(), storage.imported[1].ID)
		assert.Equal(t, "without id", storage.imported[1].Comment)
	})

	t.Run("import_v1_csv_with_and_without_header", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			body string
		}{
			{
				name: "with_header",
				body: "deliveryMethod,dose,takenAt,takenOffAt,comment\n" +
					"EV im,5,2025-01-10T12:00:00Z,,legacy\n",
			},
			{
				name: "headerless",
				body: "EV im,5,2025-01-10T12:00:00Z,,legacy\n",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				storage := &stubDoseHistoryStorage{}
				svc := NewExporterService(storage, fxtest.NewLifecycle(t), slogt.New(t))

				res, err := svc.ImportDoseHistory(ctx, strings.NewReader(tc.body), "user-1", ImportDoseHistoryOptions{
					Format: ExportCSV,
				})
				assert.NoError(t, err)
				assert.Equal(t, int64(1), res.Records)
				assert.Equal(t, uuid.Nil(), storage.imported[0].ID)
				assert.Equal(t, "EV im", storage.imported[0].DeliveryMethod)
				assert.Equal(t, float32(5), storage.imported[0].Dose)
				assert.Equal(t, "legacy", storage.imported[0].Comment)
			})
		}
	})

	t.Run("reject_estrannaise_csv_import", func(t *testing.T) {
		storage := &stubDoseHistoryStorage{}
		svc := NewExporterService(storage, fxtest.NewLifecycle(t), slogt.New(t))

		res, err := svc.ImportDoseHistory(ctx, strings.NewReader(
			"dose,days,model\n6,0,EV im\n",
		), "user-1", ImportDoseHistoryOptions{Format: ExportCSV})
		assert.Error(t, err)
		assert.Equal(t, int64(0), res.Records)
		assert.Equal(t, int64(0), res.Succeeded)
		assert.True(t, strings.Contains(err.Error(), "estrannaise"))
	})

	t.Run("import_json_with_and_without_id", func(t *testing.T) {
		storage := &stubDoseHistoryStorage{}
		svc := NewExporterService(storage, fxtest.NewLifecycle(t), slogt.New(t))

		res, err := svc.ImportDoseHistory(ctx, strings.NewReader(`[
			{"id":"019450a0-0000-7000-8000-000000000001","deliveryMethod":"EV im","dose":5,"takenAt":"2025-01-10T12:00:00Z","comment":"one"},
			{"deliveryMethod":"EV im","dose":4,"takenAt":"2025-01-14T00:00:00Z","comment":"two"}
		]`), "user-1", ImportDoseHistoryOptions{Format: ExportJSON})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), res.Records)
		assert.Equal(t, id1, storage.imported[0].ID)
		assert.Equal(t, "one", storage.imported[0].Comment)
		assert.Equal(t, uuid.Nil(), storage.imported[1].ID)
		assert.Equal(t, "two", storage.imported[1].Comment)
	})
}

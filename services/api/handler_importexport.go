package api

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/timewasted/go-accept-headers"
	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/e2clicker/services/api/openapi"
	"libdb.so/e2clicker/services/dosage"
)

type openAPIHandlerForImportExport struct {
	openapi.ServerInterface
	doseExporter *dosage.ExporterService
}

func newOpenAPIHandlerForImportExport(
	h *openAPIHandler,
	doseExporter *dosage.ExporterService,
) *openAPIHandlerForImportExport {
	return &openAPIHandlerForImportExport{
		ServerInterface: h.asHandler(),
		doseExporter:    doseExporter,
	}
}

func (h *openAPIHandlerForImportExport) asHandler() openapi.ServerInterface {
	return h
}

func (h *openAPIHandlerForImportExport) ExportDoses(w http.ResponseWriter, r *http.Request, params openapi.ExportDosesParams) {
	ctx := r.Context()
	session := sessionFromCtx(ctx)

	var format dosage.ExportFormat
acceptSearch:
	for _, t := range accept.Parse(string(params.Accept)) {
		switch t.Type + "/" + t.Subtype {
		case "text/csv":
			format = dosage.ExportCSV
			break acceptSearch
		case "application/json":
			format = dosage.ExportJSON
			break acceptSearch
		}
	}

	if format == "" {
		writeError(w, r, ErrNoAcceptableContentType, http.StatusNotAcceptable)
		return
	}

	exportExtensions, err := mime.ExtensionsByType(format.AsMIME())
	if err != nil {
		writeError(w, r, fmt.Errorf("format %q missing file extension: %w", format, err), 500)
		return
	}

	exportTime := time.Now().Format(time.RFC3339)
	exportName := fmt.Sprintf("dose-history-%s.%s", exportTime, exportExtensions[0])

	w.Header().Set("Content-Type", format.AsMIME())
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", exportName))

	n, err := h.doseExporter.ExportDoseHistory(ctx, w, session.UserSecret, dosage.ExportDoseHistoryOptions{
		Begin:  optPtr(params.Start),
		End:    optPtr(params.End),
		Format: format,
	})

	if err != nil {
		if n == 0 {
			// If we didn't write any records, we can still write the error.
			writeError(w, r, err, 0)
			return
		}
	}
}

func (h *openAPIHandlerForImportExport) ImportDoses(w http.ResponseWriter, r *http.Request, params openapi.ImportDosesParams) {
	ctx := r.Context()
	session := sessionFromCtx(ctx)

	contentType, ctParams, err := mime.ParseMediaType(string(params.ContentType))
	if err != nil {
		writeError(w, r, publicerrors.Errorf("invalid content type: %w", err), http.StatusBadRequest)
		return
	}

	if charset, ok := ctParams["charset"]; ok && charset != "utf-8" {
		writeError(w, r, publicerrors.Errorf("unsupported charset %q, UTF-8 only please", charset), http.StatusBadRequest)
		return
	}

	var format dosage.ExportFormat
	switch contentType {
	case "text/csv":
		format = dosage.ExportCSV
	case "application/json":
		format = dosage.ExportJSON
	default:
		writeError(w, r, ErrNoAcceptableContentType, http.StatusUnsupportedMediaType)
		return
	}

	result, err := h.doseExporter.ImportDoseHistory(ctx, r.Body, session.UserSecret, dosage.ImportDoseHistoryOptions{
		Format: format,
	})
	if result.Records == 0 && err != nil {
		writeError(w, r, err, 0)
		return
	}

	var oapiError *openapi.Error
	if err != nil {
		converted := convertError[errorResponse](ctx, err)
		oapiError = &converted.Body
	}

	json.NewEncoder(w).Encode(openapi.ImportDoses200JSONResponse{
		Records:   int(result.Records),
		Succeeded: int(result.Succeeded),
		Error:     oapiError,
	})
}

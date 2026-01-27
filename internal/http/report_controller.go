package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"sre/internal/usecases"
)

// NewReportController creates a report controller.
func NewReportController(s usecases.ReportService) *ReportController {
	return &ReportController{service: s}
}

type ReportController struct {
	service usecases.ReportService
}

// Routes registers report routes on r.
func (c *ReportController) Routes(r chi.Router) {
	r.Get("/report", c.getReport)
}

func (c *ReportController) getReport(w http.ResponseWriter, r *http.Request) {
	rep, err := c.service.GetReport(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "get report failed", "err", err)
		encodeError(w, "get report failed", http.StatusInternalServerError)
		return
	}
	encodeJSON(w, rep, http.StatusOK)
}

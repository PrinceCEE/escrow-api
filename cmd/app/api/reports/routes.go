package reports

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func ReportRouter(c *config.Config) chi.Router {
	h := reportHandler{c}
	r := chi.NewRouter()

	r.Post("/", h.reportTransaction)
	r.Get("/{report_id}", h.getReport)
	r.Get("/", h.getReports)

	return r
}

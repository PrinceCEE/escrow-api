package reports

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func ReportRouter(c *config.Config) chi.Router {
	rh := reportHandler{c}
	r := chi.NewRouter()

	r.Post("/", rh.reportTransaction)
	r.Get("/{report_id}", rh.getReport)
	r.Get("/", rh.getReports)

	return r
}

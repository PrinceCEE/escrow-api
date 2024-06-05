package reports

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/config"
)

func ReportRouter(c config.IConfig) chi.Router {
	h := reportHandler{c}
	r := chi.NewRouter()

	r.Post("/", h.reportTransaction)
	r.Get("/{report_id}", h.getReport)
	r.Get("/", h.getReports)

	return r
}

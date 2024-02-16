package reports

import (
	"github.com/go-chi/chi/v5"
)

func ReportRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", reportTransaction)
	r.Get("/{report_id}", getReport)
	r.Get("/", getReports)

	return r
}

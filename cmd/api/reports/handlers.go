package reports

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type reportHandler struct {
	c *config.Config
}

func (rh *reportHandler) reportTransaction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not yet implemented"))
}

func (rh *reportHandler) getReport(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not yet implemented"))
}

func (rh *reportHandler) getReports(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not yet implemented"))
}

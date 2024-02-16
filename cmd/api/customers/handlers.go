package customers

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/utils"
)

type customersHandler struct {
	c *config.Config
}

func (ch *customersHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

package businesses

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/utils"
)

type businessHandler struct {
	c *config.Config
}

func (bh *businessHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

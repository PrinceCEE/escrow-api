package customers

import (
	"net/http"

	"github.com/princecee/escrow-api/cmd/app/pkg/response"
	"github.com/princecee/escrow-api/config"
)

type customerHandler struct {
	c config.IConfig
}

func (h *customerHandler) notImplemented(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

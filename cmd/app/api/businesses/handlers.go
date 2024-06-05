package businesses

import (
	"net/http"

	"github.com/princecee/escrow-api/cmd/app/pkg/response"
	"github.com/princecee/escrow-api/config"
)

type businessHandler struct {
	c config.IConfig
}

func (h *businessHandler) notImplemented(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

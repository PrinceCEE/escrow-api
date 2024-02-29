package customers

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
)

func notImplemented(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

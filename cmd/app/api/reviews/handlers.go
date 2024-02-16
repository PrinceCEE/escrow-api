package reviews

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg"
)

func notImplemented(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

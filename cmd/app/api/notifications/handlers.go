package notifications

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
)

func getNotifications(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func getNotification(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func markAsRead(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

package notifications

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/utils"
)

type notificationHandler struct {
	c *config.Config
}

func (nh *notificationHandler) getNotifications(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (nh *notificationHandler) getNotification(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (nh *notificationHandler) markAsRead(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

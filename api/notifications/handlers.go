package notifications

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type notificationHandler struct {
	c *config.Config
}

func (nh *notificationHandler) getNotifications(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (nh *notificationHandler) getNotification(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (nh *notificationHandler) markAsRead(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

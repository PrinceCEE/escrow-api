package notifications

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func NotificationRouter(c *config.Config) chi.Router {
	nh := notificationHandler{c}
	r := chi.NewRouter()

	r.Get("/", nh.getNotifications)
	r.Get("/{notification_id}", nh.getNotification)
	r.Post("/mark-as-read", nh.markAsRead)

	return r
}

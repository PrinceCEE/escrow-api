package notifications

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func NotificationRouter(c *config.Config) chi.Router {
	h := notificationHandler{c}
	r := chi.NewRouter()

	r.Get("/", h.getNotifications)
	r.Get("/{notification_id}", h.getNotification)
	r.Post("/mark-as-read", h.markAsRead)

	return r
}

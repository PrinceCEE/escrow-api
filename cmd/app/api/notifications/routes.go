package notifications

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/config"
)

func NotificationRouter(c config.IConfig) chi.Router {
	h := notificationHandler{c}
	r := chi.NewRouter()

	r.Get("/", h.getNotifications)
	r.Get("/{notification_id}", h.getNotification)
	r.Post("/mark-as-read", h.markAsRead)

	return r
}

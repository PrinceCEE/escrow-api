package notifications

import (
	"github.com/go-chi/chi/v5"
)

func NotificationRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getNotifications)
	r.Get("/{notification_id}", getNotification)
	r.Post("/mark-as-read", markAsRead)

	return r
}

package tasks

import (
	"net/http"
	"scheduler-app/rest/middleware"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middleware.Manager) {
	mux.Handle("/api/tasks", manager.With(
		http.HandlerFunc(h.GetAllTasks),
		h.middleware.AuthenticateJWT,
	))

	mux.Handle("/api/tasks/", manager.With(
		http.HandlerFunc(h.DeleteTaskHandler),
		h.middleware.AuthenticateJWT,
	))

	mux.Handle("/api/events", manager.With(
		http.HandlerFunc(h.RealtimeHandler),
		h.middleware.AuthenticateJWT,
	))
}
package tasks

import (
	"net/http"
	"scheduler-app/rest/middleware"
)

type Router interface {
	Handle(pattern string, handler http.Handler)
}

func (h *Handler) RegisterRoutes(router Router, manager *middleware.Manager) {
	router.Handle("GET /tasks", manager.With(
		http.HandlerFunc(h.GetAllTasks),
		h.middleware.AuthenticateJWT,
	))

	router.Handle("POST /tasks", manager.With(
		http.HandlerFunc(h.CreateTask),
		h.middleware.AuthenticateJWT,
	))

	router.Handle("DELETE /tasks/{id}", manager.With(
		http.HandlerFunc(h.DeleteTaskHandler),
		h.middleware.AuthenticateJWT,
	))

	router.Handle("GET /events", manager.With(
		http.HandlerFunc(h.RealtimeHandler),
		h.middleware.AuthenticateJWT,
	))
}
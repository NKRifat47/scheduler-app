package tasks

import (
	"scheduler-app/core"
	"scheduler-app/rest/middleware"
)

type Handler struct {
	taskService   Service
	broker        *core.Broker
	middleware    *middleware.Middlewares
}

func NewHandler(taskService Service, broker *core.Broker, jwtMiddleware *middleware.Middlewares) *Handler {
	return &Handler{
		taskService: taskService,
		broker:      broker,
		middleware:  jwtMiddleware,
	}
}

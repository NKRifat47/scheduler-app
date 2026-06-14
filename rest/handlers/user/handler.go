package user

import (
	"net/http"

	"scheduler-app/config"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

type Handler struct {
	userService Service
	cnf         *config.Config
}

func NewHandler(cnf *config.Config, userService Service) *Handler {
	return &Handler{
		userService: userService,
		cnf:         cnf,
	}
}

func GetUserID(r *http.Request) int {
	val := r.Context().Value(UserIDKey)
	if val == nil {
		return 0
	}
	return val.(int)
}
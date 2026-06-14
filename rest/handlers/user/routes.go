package user

import "net/http"


type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

func (h *Handler) RegisterRoutes(router Router) {
	router.HandleFunc("POST /signup", h.SignupHandler)
	router.HandleFunc("POST /login", h.LoginHandler)
	router.HandleFunc("POST /logout", h.LogoutHandler)
}

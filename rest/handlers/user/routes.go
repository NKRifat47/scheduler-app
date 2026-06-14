package user

import "net/http"


func (h *Handler) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/api/signup", h.SignupHandler)

	mux.HandleFunc("/api/login", h.LoginHandler)

	mux.HandleFunc("/api/logout", h.LogoutHandler)
}

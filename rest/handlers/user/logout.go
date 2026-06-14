package user

import (
	"net/http"
	"scheduler-app/util"
	"time"
)

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		MaxAge:   -1,
	})

	util.SendData(w, http.StatusOK, map[string]interface{}{
		"message": "Logged out successfully",
	})
}
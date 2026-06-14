package user

import (
	"encoding/json"
	"net/http"
	"scheduler-app/util"
	"strings"
	"time"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		util.SendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := util.CreateJwt(h.cnf.JwtScretKey, util.Payload{
		Sub: user.ID,
	})
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, "Failed to generate auth token")
		return
	}

	// Still set cookie as a convenience fallback for web SSE EventSource
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	util.SendData(w, http.StatusOK, map[string]interface{}{
		"message":  "Login successful",
		"token":    token,
		"username": user.Username,
	})
}
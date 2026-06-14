package user

import (
	"encoding/json"
	"net/http"
	"scheduler-app/util"
	"strings"
	"time"
)

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
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

	if req.Username == "" || req.Password == "" {
		util.SendError(w, http.StatusBadRequest, "Username and password cannot be empty")
		return
	}

	if len(req.Password) < 6 {
		util.SendError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	userID, err := h.userService.Signup(req.Username, req.Password)
	if err != nil {
		if err.Error() == "username is already taken" {
			util.SendError(w, http.StatusConflict, "Username is already taken")
			return
		}
		util.SendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	token, err := util.CreateJwt(h.cnf.JwtScretKey, util.Payload{
		Sub: userID,
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

	util.SendData(w, http.StatusCreated, map[string]interface{}{
		"message":  "Registration successful",
		"token":    token,
		"username": req.Username,
	})
}
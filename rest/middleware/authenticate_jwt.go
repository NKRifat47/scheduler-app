package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"scheduler-app/rest/handlers/user"
	"scheduler-app/util"
)


func (m *Middlewares) AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessToken string

		// 1. Attempt to extract from Authorization header
		header := r.Header.Get("Authorization")
		if header != "" {
			headerArr := strings.Split(header, " ")
			if len(headerArr) == 2 && strings.ToLower(headerArr[0]) == "bearer" {
				accessToken = headerArr[1]
			} else if len(headerArr) == 2 {
				accessToken = headerArr[1]
			}
		}

		// 2. Fallback to query parameter (crucial for EventSource/SSE connections)
		if accessToken == "" {
			accessToken = r.URL.Query().Get("token")
		}

		// 3. Fallback to session cookie
		if accessToken == "" {
			if cookie, err := r.Cookie("session_token"); err == nil {
				accessToken = cookie.Value
			}
		}

		if accessToken == "" {
			util.SendError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		tokenParts := strings.Split(accessToken, ".")
		if len(tokenParts) != 3 {
			util.SendError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		jwtHeader := tokenParts[0]
		jwtPayload := tokenParts[1]
		jwtSignature := tokenParts[2]

		message := jwtHeader + "." + jwtPayload

		byteArrSecret := []byte(m.cnf.JwtScretKey)
		byteArrMessage := []byte(message)

		h := hmac.New(sha256.New, byteArrSecret)
		h.Write(byteArrMessage)

		hash := h.Sum(nil)
		newSignature := base64UrlEncode(hash)

		if newSignature != jwtSignature {
			util.SendError(w, http.StatusUnauthorized, "Unauthorized: invalid signature")
			return
		}

		// Decode the base64url payload to extract claims
		payloadBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(jwtPayload)
		if err != nil {
			// Try decoding with RawURLEncoding as fallback
			payloadBytes, err = base64.RawURLEncoding.DecodeString(jwtPayload)
			if err != nil {
				util.SendError(w, http.StatusUnauthorized, "Invalid token payload format")
				return
			}
		}

		var claims struct {
			Sub int   `json:"sub"`
			Exp int64 `json:"exp"`
		}
		if err := json.Unmarshal(payloadBytes, &claims); err != nil {
			util.SendError(w, http.StatusUnauthorized, "Invalid token claims format")
			return
		}

		// Check if token has expired
		if claims.Exp > 0 && claims.Exp < time.Now().Unix() {
			util.SendError(w, http.StatusUnauthorized, "Token has expired")
			return
		}

		// Inject the user ID into the request context using the shared context key
		ctx := context.WithValue(r.Context(), user.UserIDKey, claims.Sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func base64UrlEncode(data []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

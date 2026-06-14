package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scheduler-app/core"
	"scheduler-app/domain"
	"scheduler-app/rest/handlers/user"
	"scheduler-app/util"
	"time"
)

func (h *Handler) RealtimeHandler(w http.ResponseWriter, r *http.Request) {
	userID := user.GetUserID(r)
	if userID == 0 {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		util.SendError(w, http.StatusInternalServerError, "Streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := make(chan *domain.Task, 10)
	h.broker.Register <- core.ClientConn{Ch: ch, UserID: userID}

	defer func() {
		h.broker.Unregister <- ch
	}()

	flusher.Flush()

	pingTicker := time.NewTicker(15 * time.Second)
	defer pingTicker.Stop()

	log.Printf("SSE connection established for User ID %d", userID)

	for {
		select {
		case task, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(task)
			if err != nil {
				continue
			}
			_, _ = fmt.Fprintf(w, "event: task_triggered\ndata: %s\n\n", data)
			flusher.Flush()

		case <-pingTicker.C:
			// Send periodic comment line to keep the connection active
			_, _ = fmt.Fprintf(w, ": ping\n\n")
			flusher.Flush()

		case <-r.Context().Done():
			log.Printf("SSE connection closed by context for User ID %d", userID)
			return
		}
	}
}
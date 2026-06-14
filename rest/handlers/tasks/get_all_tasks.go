package tasks

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"scheduler-app/domain"
	"scheduler-app/rest/handlers/user"
	"scheduler-app/util"
)

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userID := user.GetUserID(r)
	if userID == 0 {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	switch r.Method {
	case http.MethodGet:
		tasks, err := h.taskService.GetTasksForUser(userID)
		if err != nil {
			util.SendError(w, http.StatusInternalServerError, "Failed to fetch tasks")
			return
		}
		if tasks == nil {
			tasks = []*domain.Task{}
		}
		util.SendData(w, http.StatusOK, tasks)

	case http.MethodPost:
		var req struct {
			Title         string `json:"title"`
			Description   string `json:"description"`
			ScheduledTime string `json:"scheduled_time"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.SendError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		req.Title = strings.TrimSpace(req.Title)
		req.Description = strings.TrimSpace(req.Description)

		scheduled, err := time.Parse(time.RFC3339, req.ScheduledTime)
		if err != nil {
			// Try a backup parser if browser didn't send exact RFC3339
			scheduled, err = time.Parse("2006-01-02T15:04", req.ScheduledTime)
			if err != nil {
				util.SendError(w, http.StatusBadRequest, "Invalid scheduled time format. Use ISO/RFC3339")
				return
			}
			// Use local time zone since datepicker input has no offset
			localTime := time.Date(
				scheduled.Year(), scheduled.Month(), scheduled.Day(),
				scheduled.Hour(), scheduled.Minute(), scheduled.Second(),
				0, time.Local,
			)
			scheduled = localTime
		}

		task, err := h.taskService.CreateTask(userID, req.Title, req.Description, scheduled)
		if err != nil {
			if err.Error() == "title is required" || err.Error() == "scheduled time must be in the future" {
				util.SendError(w, http.StatusBadRequest, err.Error())
				return
			}
			log.Printf("Error creating task: %v", err)
			util.SendError(w, http.StatusInternalServerError, "Failed to save task")
			return
		}

		util.SendData(w, http.StatusCreated, task)

	default:
		util.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
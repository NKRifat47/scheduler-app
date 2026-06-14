package tasks

import (
	"errors"
	"net/http"
	"scheduler-app/rest/handlers/user"
	"scheduler-app/util"
	"strconv"
	"strings"
)

// DeleteTaskHandler deletes a specific task
func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		util.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := user.GetUserID(r)
	if userID == 0 {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		util.SendError(w, http.StatusBadRequest, "Missing task ID")
		return
	}

	taskID, err := strconv.Atoi(parts[3])
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = h.taskService.DeleteTask(userID, taskID)
	if err != nil {
		if errors.New("task not found or does not belong to user").Error() == err.Error() {
			util.SendError(w, http.StatusNotFound, err.Error())
			return
		}
		util.SendError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}
	util.SendData(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

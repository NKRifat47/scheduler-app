package tasks

import (
	"errors"
	"net/http"
	"scheduler-app/rest/handlers/user"
	"scheduler-app/util"
	"strconv"
)

// DeleteTaskHandler deletes a specific task
func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := user.GetUserID(r)
	if userID == 0 {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	taskIDStr := r.PathValue("id")
	if taskIDStr == "" {
		util.SendError(w, http.StatusBadRequest, "Missing task ID")
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
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

package tasks

import (
	"scheduler-app/domain"
	"time"
)

type Service interface {
	CreateTask(userID int, title, description string, scheduledTime time.Time) (*domain.Task, error)
	GetTasksForUser(userID int) ([]*domain.Task, error)
	DeleteTask(userID, taskID int) error
}

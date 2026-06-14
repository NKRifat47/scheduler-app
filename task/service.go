package task

import (
	"errors"
	"scheduler-app/domain"
	"time"
)

type service struct {
	taskRepo TaskRepo
}

func NewService(taskRepo TaskRepo) Service {
	return &service{
		taskRepo: taskRepo,
	}
}

func (s *service) CreateTask(userID int, title, description string, scheduledTime time.Time) (*domain.Task, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if scheduledTime.Before(time.Now()) {
		return nil, errors.New("scheduled time must be in the future")
	}
	return s.taskRepo.CreateTask(userID, title, description, scheduledTime)
}

func (s *service) GetTasksForUser(userID int) ([]*domain.Task, error) {
	return s.taskRepo.GetTasksForUser(userID)
}

func (s *service) DeleteTask(userID, taskID int) error {
	return s.taskRepo.DeleteTask(userID, taskID)
}

func (s *service) GetPendingTasks() ([]*domain.Task, error) {
	return s.taskRepo.GetPendingTasks()
}

func (s *service) MarkTaskAsTriggered(taskID int) error {
	return s.taskRepo.MarkTaskAsTriggered(taskID)
}

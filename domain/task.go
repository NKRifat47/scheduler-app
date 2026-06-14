package domain

import "time"

type Task struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ScheduledTime time.Time `json:"scheduled_time"`
	Triggered     bool      `json:"triggered"`
}
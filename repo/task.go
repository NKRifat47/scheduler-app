package repo

import (
	"errors"
	"scheduler-app/domain"
	"scheduler-app/task"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskRepo interface {
	task.TaskRepo
}

type taskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) TaskRepo {
	return &taskRepo{
		db: db,
	}
}

func (r *taskRepo) CreateTask(userID int, title, description string, scheduledTime time.Time) (*domain.Task, error) {
	var t domain.Task
	query := `
	INSERT INTO tasks (user_id, title, description, scheduled_time)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, title, description, scheduled_time, triggered`

	err := r.db.QueryRow(query, userID, title, description, scheduledTime).
		Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.ScheduledTime, &t.Triggered)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepo) GetTasksForUser(userID int) ([]*domain.Task, error) {
	query := `
	SELECT id, user_id, title, description, scheduled_time, triggered 
	FROM tasks 
	WHERE user_id = $1 
	ORDER BY scheduled_time ASC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.ScheduledTime, &t.Triggered)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepo) DeleteTask(userID, taskID int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, taskID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("task not found or does not belong to user")
	}
	return nil
}

func (r *taskRepo) GetPendingTasks() ([]*domain.Task, error) {
	query := `
	SELECT id, user_id, title, description, scheduled_time, triggered 
	FROM tasks 
	WHERE triggered = FALSE AND scheduled_time <= $1`

	rows, err := r.db.Query(query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.ScheduledTime, &t.Triggered)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepo) MarkTaskAsTriggered(taskID int) error {
	query := `UPDATE tasks SET triggered = TRUE WHERE id = $1`
	_, err := r.db.Exec(query, taskID)
	return err
}

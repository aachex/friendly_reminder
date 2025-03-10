package repository

import (
	"database/sql"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

type TasksRepository interface {
	// AddItem добавляет новую задачу в список пользователя.
	AddTask(value, userEmail string) error

	// GetList возвращает список дел пользователя с указанным email.
	GetList(userEmail string) ([]models.Task, error)
}

type tasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) TasksRepository {
	return &tasksRepository{
		db: db,
	}
}

func (r *tasksRepository) AddTask(value, userEmail string) error {
	_, err := r.db.Exec(`
		INSERT INTO tasks(value, number_in_list, user_email) VALUES(
			$1,
			(SELECT COUNT(*) FROM tasks WHERE user_email = $2) + 1,
			$2)`,
		value, userEmail)
	return err
}

// GetList возвращает список дел пользователя с указанным email.
func (r *tasksRepository) GetList(userEmail string) ([]models.Task, error) {
	rows, err := r.db.Query("SELECT * FROM tasks WHERE user_email = $1 ORDER BY number_in_list", userEmail)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	var task models.Task
	for rows.Next() {
		rows.Scan(&task.Id, &task.UserEmail, &task.Value, &task.NumberInList)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

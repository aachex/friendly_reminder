package repository

import (
	"database/sql"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

type TasksRepository interface {
	// AddItem добавляет новую задачу в список пользователя. Возвращает id созданноё задачи.
	AddTask(value, userEmail string) (int64, error)

	// DeleteTask удаляет задачу по указанному id.
	DeleteTask(id int64) error

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

// AddItem добавляет новую задачу в список пользователя. Возвращает id созданноё задачи.
func (r *tasksRepository) AddTask(value, userEmail string) (int64, error) {
	res, err := r.db.Exec(`INSERT INTO tasks(value, user_email) VALUES($1, $2)`, value, userEmail)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteTask удаляет задачу по указанному id.
func (r *tasksRepository) DeleteTask(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE task_id = $1`, id)
	return err
}

// GetList возвращает список дел пользователя с указанным email.
func (r *tasksRepository) GetList(userEmail string) ([]models.Task, error) {
	rows, err := r.db.Query("SELECT * FROM tasks WHERE user_email = $1", userEmail)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	var task models.Task
	for rows.Next() {
		rows.Scan(&task.Id, &task.UserEmail, &task.Value)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

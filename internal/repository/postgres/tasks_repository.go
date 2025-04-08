package sqlite

import (
	"context"
	"database/sql"
	"sync"

	"github.com/artemwebber1/friendly_reminder/internal/models"
)

type TasksRepository struct {
	mu sync.Mutex
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{
		db: db,
		mu: sync.Mutex{},
	}
}

// AddItem добавляет новую задачу в список пользователя. Возвращает id созданной задачи.
func (r *TasksRepository) AddTask(ctx context.Context, value, userEmail string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	row := r.db.QueryRowContext(ctx, "INSERT INTO tasks(value, user_email) VALUES($1, $2) RETURNING task_id", value, userEmail)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteTask удаляет задачу по указанному id.
func (r *TasksRepository) DeleteTask(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE task_id = $1", id)
	return err
}

// GetList возвращает список дел пользователя с указанным email.
func (r *TasksRepository) GetList(ctx context.Context, userEmail string) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM tasks WHERE user_email = $1", userEmail)
	if err != nil {
		return []models.Task{}, err
	}

	tasks := make([]models.Task, 0)
	var task models.Task
	for rows.Next() {
		rows.Scan(&task.Id, &task.UserEmail, &task.Value)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// ClearList очищает список указанного пользователя.
func (r *TasksRepository) ClearList(ctx context.Context, userEmail string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE user_email = $1", userEmail)
	return err
}

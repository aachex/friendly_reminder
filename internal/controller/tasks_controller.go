package controller

import (
	"net/http"
	"strconv"

	"github.com/artemwebber1/friendly_reminder/internal/repository"
	mw "github.com/artemwebber1/friendly_reminder/pkg/middleware"
)

type TasksController struct {
	tasksRepo repository.TasksRepository
}

func NewTasksController(tr repository.TasksRepository) *TasksController {
	return &TasksController{
		tasksRepo: tr,
	}
}

func (c *TasksController) AddEndpoints(mux *http.ServeMux) {
	mux.HandleFunc(
		"POST /new-task",
		mw.UseLogging(mw.UseAuthorization(c.CreateTask)))

	mux.HandleFunc(
		"DELETE /del-task",
		mw.UseLogging(mw.UseAuthorization(c.DeleteTask)))
}

// CreateTask создаёт новую задачу в списке пользователя.
//
// Обрабатывает POST запросы по пути '/new-task'.
func (c *TasksController) CreateTask(w http.ResponseWriter, r *http.Request) {
	rawJwt := getRawJwtFromHeader(r.Header)
	jwtClaims, err := readJWT(rawJwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	type newTask struct {
		Value string `json:"value"`
	}

	task, err := readBody[newTask](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.tasksRepo.AddTask(task.Value, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task was created succesfully"))
}

// DeleteTask удаляет задачу из списка пользователя.
//
// Обрабатывает DELETE запросы по пути '/del-task'.
func (c *TasksController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.tasksRepo.DeleteTask(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

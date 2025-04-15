package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/artemwebber1/friendly_reminder/internal/config"
	"github.com/artemwebber1/friendly_reminder/internal/logging"
	"github.com/artemwebber1/friendly_reminder/internal/models"
	"github.com/artemwebber1/friendly_reminder/pkg/authorization"
)

type tasksRepository interface {
	// AddItem добавляет новую задачу в список пользователя. Возвращает id созданной задачи.
	AddTask(ctx context.Context, value, userEmail string) (int64, error)

	// DeleteTask удаляет задачу по указанному id.
	DeleteTask(ctx context.Context, id int64) error

	// GetList возвращает список дел пользователя с указанным email.
	GetList(ctx context.Context, userEmail string) ([]models.Task, error)

	// ClearList очищает список указанного пользователя.
	ClearList(ctx context.Context, userEmail string) error
}

type TasksController struct {
	tasksRepo tasksRepository
	usersRepo usersRepository
	cfg       *config.Config
}

func NewTasksController(tr tasksRepository, ur usersRepository, cfg *config.Config) *TasksController {
	return &TasksController{
		tasksRepo: tr,
		usersRepo: ur,
		cfg:       cfg,
	}
}

func (c *TasksController) AddEndpoints(mux *http.ServeMux) {
	mux.HandleFunc(
		"POST "+c.cfg.Prefix+"/tasks/new",
		logging.Middleware(authorization.Middleware(c.CreateTask)),
	)

	mux.HandleFunc(
		"GET "+c.cfg.Prefix+"/tasks/list",
		logging.Middleware(authorization.Middleware(c.GetList)),
	)

	mux.HandleFunc(
		"DELETE "+c.cfg.Prefix+"/tasks/clear-list",
		logging.Middleware(authorization.Middleware(c.ClearList)),
	)

	mux.HandleFunc(
		"DELETE "+c.cfg.Prefix+"/tasks/del",
		logging.Middleware(authorization.Middleware(c.DeleteTask)),
	)
}

// CreateTask создаёт новую задачу в списке пользователя.
//
// Обрабатывает POST запросы по пути '/tasks/new'.
func (c *TasksController) CreateTask(w http.ResponseWriter, r *http.Request) {
	rawJwt := authorization.FromHeader(r.Header)
	jwtClaims, err := authorization.GetClaims(rawJwt, jwtKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if !c.usersRepo.EmailExists(r.Context(), email) {
		http.Error(w, errInvalidInvalidEmail.Error(), http.StatusForbidden)
		return
	}

	type newTask struct {
		Id    int64  `json:"task_id"`
		Value string `json:"value"`
	}

	task, err := readBody[newTask](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := c.tasksRepo.AddTask(r.Context(), task.Value, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.Id = id

	w.WriteHeader(http.StatusCreated)
	writeJson(w, task)
}

// GetList получает список пользователя.
//
// Обрабатывает GET запросы по пути '/tasks/list'.
func (c *TasksController) GetList(w http.ResponseWriter, r *http.Request) {
	rawJwt := authorization.FromHeader(r.Header)
	jwtClaims, err := authorization.GetClaims(rawJwt, jwtKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if !c.usersRepo.EmailExists(r.Context(), email) {
		http.Error(w, errInvalidInvalidEmail.Error(), http.StatusForbidden)
		return
	}

	list, err := c.tasksRepo.GetList(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, &list)
}

// ClearList удаляет все задачи из списка пользователя.
//
// Обрабатывает DELETE запросы по пути '/tasks/clear-list'.
func (c *TasksController) ClearList(w http.ResponseWriter, r *http.Request) {
	rawJwt := authorization.FromHeader(r.Header)
	jwtClaims, err := authorization.GetClaims(rawJwt, jwtKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if !c.usersRepo.EmailExists(r.Context(), email) {
		http.Error(w, errInvalidInvalidEmail.Error(), http.StatusForbidden)
		return
	}

	err = c.tasksRepo.ClearList(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteTask удаляет задачу из списка пользователя.
//
// Обрабатывает DELETE запросы по пути '/tasks/del'.
func (c *TasksController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	rawJwt := authorization.FromHeader(r.Header)
	jwtClaims, err := authorization.GetClaims(rawJwt, jwtKey())
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	email, err := jwtClaims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if !c.usersRepo.EmailExists(r.Context(), email) {
		http.Error(w, errInvalidInvalidEmail.Error(), http.StatusForbidden)
		return
	}

	taskId, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.tasksRepo.DeleteTask(r.Context(), taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

package controller

import "net/http"

type usersController struct {
}

// GetTodoList возвращает список пользователя.
//
// Обрабатывает GET запросы по пути '/get-todo'.
//
// Требует jwt. При отсутствии возвращает ошибку 403.
func (c *usersController) GetTodoList(w http.ResponseWriter, r *http.Request)

package models

// Task - это задача, входящая в список пользователя.
type Task struct {
	Id        int64  `json:"task_id"`
	UserEmail string `json:"user_email"`
	Value     string `json:"value"`
}

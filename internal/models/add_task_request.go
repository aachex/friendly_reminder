package models

type AddTaskRequest struct {
	TaskValue string `json:"task_value"`
	UserEmail string `json:"email"`
}

package model

import "time"

// Task representa a estrutura de uma tarefa no sistema

type Task struct {
	ID          string    `json:"id" example:"5386b9cf-e4c5-40e6-bc29-fc901dc04290"`
	Title       string    `json:"title" example:"Estudar Golang"`
	Description string    `json:"description" example:"Aprender documentação com Swagger"`
	Status      string    `json:"status" example:"pending"`
	Priority    string    `json:"priority" example:"high"`
	DueDate     time.Time `json:"due_date" example:"2026-04-10T00:00:00Z"`
	CreatedAt   time.Time `json:"created_at" example:"2026-03-10T10:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-03-10T10:00:00Z"`
}

type CreateTaskRequest struct { //request do POST
	Title       string    `json:"title" example:"Estudar Golang"`
	Description string    `json:"description" example:"Aprender documentação com Swagger"`
	Priority    string    `json:"priority" example:"high" enums:"low,medium,high"`
	DueDate     time.Time `json:"due_date" example:"2026-04-10T00:00:00Z"`
}

type UpdateTaskRequest struct { //request PUT
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date" example:"2026-04-10T00:00:00Z"`
}

package model

import "time"

// Task representa a estrutura de uma tarefa no sistema
type Task struct {
	ID          string    `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Status      string    `json:"status" bson:"status"`
	Priority    string    `json:"priority" bson:"priority"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

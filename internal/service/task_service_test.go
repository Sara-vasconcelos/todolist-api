package service

import (
	"testing"
	"time"

	"todolist-api/internal/model"
)

type MockTaskRepository struct {
	status string
}

func (m *MockTaskRepository) CreateTask(task *model.Task) error {
	return nil
}

func (m *MockTaskRepository) GetTasks(filter map[string]interface{}) ([]model.Task, error) {
	return []model.Task{}, nil
}

func (m *MockTaskRepository) GetTaskByID(id string) (*model.Task, error) {
	return &model.Task{
		ID:       id,
		Title:    "Teste",
		Status:   m.status,
		Priority: "medium",
		DueDate:  time.Now().Add(24 * time.Hour),
	}, nil
}

func (m *MockTaskRepository) UpdateTask(id string, task *model.Task) error {
	return nil
}

func (m *MockTaskRepository) DeleteTask(id string) error {
	return nil
}

func TestCreateTaskSuccess(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	task := &model.Task{
		Title:    "Nova Task",
		Status:   "pending",
		Priority: "high",
		DueDate:  time.Now().Add(24 * time.Hour),
	}

	err := service.CreateTask(task)

	if err != nil {
		t.Errorf("esperava sucesso, mas recebeu erro: %v", err)
	}
}

func TestCreateTaskInvalidTitle(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	task := &model.Task{
		Title:    "a",
		Status:   "pending",
		Priority: "high",
		DueDate:  time.Now().Add(24 * time.Hour),
	}

	err := service.CreateTask(task)

	if err == nil {
		t.Error("esperava erro para título inválido")
	}
}

func TestCreateTaskPastDueDate(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	task := &model.Task{
		Title:    "Task Teste",
		Status:   "pending",
		Priority: "high",
		DueDate:  time.Now().Add(-24 * time.Hour),
	}

	err := service.CreateTask(task)

	if err == nil {
		t.Error("esperava erro para data no passado")
	}
}

func TestCreateTaskInvalidStatus(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	task := &model.Task{
		Title:    "Task Teste",
		Status:   "invalid_status",
		Priority: "high",
		DueDate:  time.Now().Add(24 * time.Hour),
	}

	err := service.CreateTask(task)

	if err == nil {
		t.Error("esperava erro para status inválido")
	}
}

func TestCreateTaskInvalidPriority(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	task := &model.Task{
		Title:    "Task Teste",
		Status:   "pending",
		Priority: "invalid",
		DueDate:  time.Now().Add(24 * time.Hour),
	}

	err := service.CreateTask(task)

	if err == nil {
		t.Error("esperava erro para prioridade inválida")
	}
}

func TestListTasks(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	filter := map[string]interface{}{}

	tasks, err := service.ListTasks(filter)

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}

	if tasks == nil {
		t.Error("esperava lista de tasks")
	}
}

func TestGetTask(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	task, err := service.GetTask("123")

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}

	if task == nil {
		t.Error("esperava task retornada")
	}
}

func TestUpdateTaskSuccess(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	task := &model.Task{
		Title:    "Updated Task",
		Status:   "pending",
		Priority: "medium",
		DueDate:  time.Now().Add(24 * time.Hour),
	}

	err := service.UpdateTask("123", task)

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}

func TestUpdateTaskCompleted(t *testing.T) {

	repo := &MockTaskRepository{status: "completed"}
	service := NewTaskService(repo)

	task := &model.Task{
		Title:    "Updated Task",
		Status:   "pending",
		Priority: "medium",
		DueDate:  time.Now().Add(24 * time.Hour),
	}

	err := service.UpdateTask("123", task)

	if err == nil {
		t.Error("esperava erro ao editar task completed")
	}
}

func TestDeleteTask(t *testing.T) {

	repo := &MockTaskRepository{status: "pending"}
	service := NewTaskService(repo)

	err := service.DeleteTask("123")

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}
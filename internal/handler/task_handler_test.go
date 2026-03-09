package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todolist-api/internal/model"
	"todolist-api/internal/service"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

// Mock do repository
type MockRepository struct{}

func (m *MockRepository) CreateTask(task *model.Task) error {
	return nil
}

func (m *MockRepository) GetTasks(filter map[string]interface{}) ([]model.Task, error) {
	return []model.Task{}, nil
}

func (m *MockRepository) GetTaskByID(id string) (*model.Task, error) {
	return &model.Task{
		ID:       id,
		Title:    "Teste",
		Status:   "pending",
		Priority: "high",
		DueDate:  time.Now().Add(24 * time.Hour),
	}, nil
}

func (m *MockRepository) UpdateTask(id string, task *model.Task) error {
	return nil
}

func (m *MockRepository) DeleteTask(id string) error {
	return nil
}

func TestCreateTaskHandler(t *testing.T) {

	repo := &MockRepository{}
	svc := service.NewTaskService(repo)

	handler := NewTaskHandler(svc)

	body := `{
		"title":"Test Task",
		"status":"pending",
		"priority":"high",
		"due_date":"2030-12-31T00:00:00Z"
	}`

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperava status 201, recebeu %d", resp.StatusCode)
	}
}

func TestCreateTaskInvalidJSON(t *testing.T) {

	repo := &MockRepository{}
	svc := service.NewTaskService(repo)
	handler := NewTaskHandler(svc)

	body := `{"title":`

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperava status 400, recebeu %d", resp.StatusCode)
	}
}

func TestListTasks(t *testing.T) {

	repo := &MockRepository{}
	svc := service.NewTaskService(repo)
	handler := NewTaskHandler(svc)

	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	handler.ListTasks(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
	}
}

func TestGetTask(t *testing.T) {

	repo := &MockRepository{}
	svc := service.NewTaskService(repo)
	handler := NewTaskHandler(svc)

	id := uuid.New().String()

	req := httptest.NewRequest("GET", "/tasks/"+id, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{
		"id": id,
	})

	handler.GetTask(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
	}
}

func TestUpdateTask(t *testing.T) {

	repo := &MockRepository{}
	svc := service.NewTaskService(repo)
	handler := NewTaskHandler(svc)

	id := uuid.New().String()

	body := `{
		"title":"Updated Task",
		"status":"pending",
		"priority":"medium",
		"due_date":"2030-12-31T00:00:00Z"
	}`

	req := httptest.NewRequest("PUT", "/tasks/"+id, bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{
		"id": id,
	})

	handler.UpdateTask(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperava status 200, recebeu %d", resp.StatusCode)
	}
}

func TestDeleteTask(t *testing.T) {

	repo := &MockRepository{}
	svc := service.NewTaskService(repo)
	handler := NewTaskHandler(svc)

	id := uuid.New().String()

	req := httptest.NewRequest("DELETE", "/tasks/"+id, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{
		"id": id,
	})

	handler.DeleteTask(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("esperava status 204, recebeu %d", resp.StatusCode)
	}
}
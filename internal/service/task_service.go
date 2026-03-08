package service

import (
	"errors"
	"time"

	"todolist-api/internal/repository"
)

// TaskService é responsável pelas regras de negócio
type TaskService struct {
	repo *repository.TaskRepository
}

// constantes para validar status e prioridade
var validStatus = map[string]bool{
	"pending":     true,
	"in_progress": true,
	"completed":   true,
	"cancelled":   true,
}

var validPriority = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
}

// NewTaskService cria uma instância do service
func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// ------------------- CREATE -------------------
func (s *TaskService) CreateTask(task *repository.Task) error {
	// Validações
	if len(task.Title) < 3 || len(task.Title) > 100 {
		return errors.New("título deve ter entre 3 e 100 caracteres")
	}

	if _, ok := validStatus[task.Status]; !ok {
		return errors.New("status inválido")
	}

	if _, ok := validPriority[task.Priority]; !ok {
		return errors.New("prioridade inválida")
	}

	if task.DueDate.Before(time.Now()) {
		return errors.New("data de vencimento não pode ser no passado")
	}

	// Se passou nas validações, chama o repository
	return s.repo.CreateTask(task)
}

// ------------------- LIST -------------------
func (s *TaskService) ListTasks(filter map[string]interface{}) ([]repository.Task, error) {
	return s.repo.GetTasks(filter)
}

// ------------------- GET BY ID -------------------
func (s *TaskService) GetTask(id string) (*repository.Task, error) {
	return s.repo.GetTaskByID(id)
}

// ------------------- UPDATE -------------------
func (s *TaskService) UpdateTask(id string, updatedTask *repository.Task) error {
	// Busca a task existente
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		return err
	}

	// Regra: não edita completed
	if task.Status == "completed" {
		return errors.New("tarefas com status completed não podem ser editadas")
	}

	// Validações
	if len(updatedTask.Title) < 3 || len(updatedTask.Title) > 100 {
		return errors.New("título deve ter entre 3 e 100 caracteres")
	}

	if _, ok := validStatus[updatedTask.Status]; !ok {
		return errors.New("status inválido")
	}

	if _, ok := validPriority[updatedTask.Priority]; !ok {
		return errors.New("prioridade inválida")
	}

	if updatedTask.DueDate.Before(time.Now()) {
		return errors.New("data de vencimento não pode ser no passado")
	}

	return s.repo.UpdateTask(id, updatedTask)
}

// ------------------- DELETE -------------------
func (s *TaskService) DeleteTask(id string) error {
	return s.repo.DeleteTask(id)
}

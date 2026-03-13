package service

import (
	"errors"
	"time"

	"todolist-api/internal/model"

	log "github.com/sirupsen/logrus"
)

// Interface que define os métodos que o repository precisa implementar
type TaskRepository interface {
	CreateTask(task *model.Task) error
	GetTasks(filter map[string]interface{}) ([]model.Task, error)
	GetTaskByID(id string) (*model.Task, error)
	UpdateTask(id string, task *model.Task) error
	DeleteTask(id string) error
}

// TaskService é responsável pelas regras de negócio
type TaskService struct {
	repo TaskRepository
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
func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// ------------------- CREATE -------------------
func (s *TaskService) CreateTask(task *model.Task) error {

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "CreateTask",
		"title":   task.Title,
	}).Info("iniciando criação de task")

	// Validações
	if len(task.Title) < 3 || len(task.Title) > 100 {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "CreateTask",
			"title":   task.Title,
		}).Warn("título inválido")

		return errors.New("título deve ter entre 3 e 100 caracteres")
	}

	if _, ok := validStatus[task.Status]; !ok {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "CreateTask",
			"status":  task.Status,
		}).Warn("status inválido")

		return errors.New("status inválido")
	}

	if _, ok := validPriority[task.Priority]; !ok {

		log.WithFields(log.Fields{
			"service":  "TaskService",
			"method":   "CreateTask",
			"priority": task.Priority,
		}).Warn("prioridade inválida")

		return errors.New("prioridade inválida")
	}

	if task.DueDate.Before(time.Now()) {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "CreateTask",
			"dueDate": task.DueDate,
		}).Warn("data de vencimento inválida")

		return errors.New("data de vencimento não pode ser no passado")
	}

	// Se passou nas validações, chama o repository
	err := s.repo.CreateTask(task)

	if err != nil {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "CreateTask",
			"title":   task.Title,
			"error":   err.Error(),
		}).Error("erro ao criar task no repository")

		return err
	}

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "CreateTask",
		"task_id": task.ID,
	}).Info("task criada com sucesso")

	return nil
}

// ------------------- LIST -------------------
func (s *TaskService) ListTasks(filter map[string]interface{}) ([]model.Task, error) {

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "ListTasks",
		"filters": filter,
	}).Info("listando tasks")

	tasks, err := s.repo.GetTasks(filter)

	if err != nil {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "ListTasks",
			"error":   err.Error(),
		}).Error("erro ao buscar tasks")

		return nil, err
	}

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "ListTasks",
		"total":   len(tasks),
	}).Info("tasks retornadas com sucesso")

	return tasks, nil
}

// ------------------- GET BY ID -------------------
func (s *TaskService) GetTask(id string) (*model.Task, error) {

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "GetTask",
		"task_id": id,
	}).Info("buscando task por id")

	task, err := s.repo.GetTaskByID(id)

	if err != nil {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "GetTask",
			"task_id": id,
			"error":   err.Error(),
		}).Error("erro ao buscar task")

		return nil, err
	}

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "GetTask",
		"task_id": id,
	}).Info("task encontrada")

	return task, nil
}

// ------------------- UPDATE -------------------
func (s *TaskService) UpdateTask(id string, updatedTask *model.Task) error {

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "UpdateTask",
		"task_id": id,
	}).Info("iniciando atualização da task")

	// Busca a task existente
	task, err := s.repo.GetTaskByID(id)
	if err != nil {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "UpdateTask",
			"task_id": id,
			"error":   err.Error(),
		}).Error("erro ao buscar task para atualização")

		return err
	}

	// Regra: não edita completed
	if task.Status == "completed" {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "UpdateTask",
			"task_id": id,
		}).Warn("tentativa de editar task completed")

		return errors.New("tarefas com status completed não podem ser editadas")
	}

	// Validações
	if len(updatedTask.Title) < 3 || len(updatedTask.Title) > 100 {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "UpdateTask",
		}).Warn("título inválido")

		return errors.New("título deve ter entre 3 e 100 caracteres")
	}

	if _, ok := validStatus[updatedTask.Status]; !ok {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "UpdateTask",
			"status":  updatedTask.Status,
		}).Warn("status inválido")

		return errors.New("status inválido")
	}

	if _, ok := validPriority[updatedTask.Priority]; !ok {

		log.WithFields(log.Fields{
			"service":  "TaskService",
			"method":   "UpdateTask",
			"priority": updatedTask.Priority,
		}).Warn("prioridade inválida")

		return errors.New("prioridade inválida")
	}

	if updatedTask.DueDate.Before(time.Now()) {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "UpdateTask",
		}).Warn("data de vencimento inválida")

		return errors.New("data de vencimento não pode ser no passado")
	}

	err = s.repo.UpdateTask(id, updatedTask)

	if err != nil {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "UpdateTask",
			"task_id": id,
			"error":   err.Error(),
		}).Error("erro ao atualizar task")

		return err
	}

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "UpdateTask",
		"task_id": id,
	}).Info("task atualizada com sucesso")

	return nil
}

// ------------------- DELETE -------------------
func (s *TaskService) DeleteTask(id string) error {

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "DeleteTask",
		"task_id": id,
	}).Info("removendo task")

	err := s.repo.DeleteTask(id)

	if err != nil {

		log.WithFields(log.Fields{
			"service": "TaskService",
			"method":  "DeleteTask",
			"task_id": id,
			"error":   err.Error(),
		}).Error("erro ao deletar task")

		return err
	}

	log.WithFields(log.Fields{
		"service": "TaskService",
		"method":  "DeleteTask",
		"task_id": id,
	}).Info("task deletada com sucesso")

	return nil
}

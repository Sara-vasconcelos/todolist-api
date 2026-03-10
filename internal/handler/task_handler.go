package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"todolist-api/internal/model"
	"todolist-api/internal/repository"
	"todolist-api/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// TaskHandler lida com requisições HTTP relacionadas a tarefas
type TaskHandler struct {
	service *service.TaskService
}

// NewTaskHandler cria uma instância do handler
func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// ------------------- CREATE TASK -------------------
/*
w http.ResponseWriter : escreve a resposta HTTP
r *http.Request: contém os dados da requisição (body, headers, query params).
*/

// CreateTask cria uma nova tarefa
// @Summary Criar tarefa
// @Description Cria uma nova tarefa.
// @Description
// @Description Regras:
// @Description - Título é obrigatório (mínimo 3 e máximo 100 caracteres)
// @Description - Status permitidos: pending, in_progress, completed, cancelled
// @Description - Prioridade: low, medium, high
// @Description - Data de vencimento não pode estar no passado
// @Tags tarefas
// @Accept json
// @Produce json
// @Param task body model.CreateTaskRequest true "Dados da nova tarefa"
// @Success 201 {object} model.Task
// @Failure 400 {object} model.ErrorResponse "Título inválido (3-100 caracteres)"
// @Failure 400 {object} model.ErrorResponse "Prioridade inválida"
// @Failure 400 {object} model.ErrorResponse "Data de vencimento no passado"
// @Failure 500 {object} model.ErrorResponse "Erro interno"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

	var req model.CreateTaskRequest

	/*Pega o JSON enviado e converte em uma struct do tipo CreateTaskRequest*/
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(log.Fields{
			"handler": "CreateTask",
			"error":   err.Error(),
		}).Error("erro ao decodificar JSON da requisição")

		http.Error(w, "dados inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	task := model.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	// Define status padrão se não tiver sido fornecido
	if task.Status == "" {
		task.Status = "pending"
	}

	/*Chama o TaskService para criar a task, passando um ponteiro para task
	Valida as regras de negocio antes de salvar no banco*/
	if err := h.service.CreateTask(&task); err != nil {
		log.WithFields(log.Fields{
			"handler": "CreateTask",
			"title":   task.Title,
			"error":   err.Error(),
		}).Error("erro ao criar tarefa")

		http.Error(w, err.Error(), http.StatusBadRequest) //retorna 400 em caso de erro
		return
	}

	log.WithFields(log.Fields{
		"handler": "CreateTask",
		"task_id": task.ID,
		"status":  task.Status,
	}).Info("tarefa criada com sucesso")

	w.WriteHeader(http.StatusCreated) //retorna 201  em caso de sucesso
	json.NewEncoder(w).Encode(task)   //envia o JSON da task criada de volta para o cliente, incluindo id, created_at e updated_at.
}

// ------------------- LIST TASKS -------------------

// @Summary Listar todas as tarefas
// @Description Retorna todas as tarefas, com filtros opcionais de status e prioridade
// @Tags tarefas
// @Accept json
// @Produce json
// @Param status query string false "Filtrar tarefas por status" Enums(pending,in_progress,completed,cancelled)
// @Param priority query string false "Filtrar tarefas por prioridade" Enums(low,medium,high)
// @Success 200 {array} model.Task
// @Failure 500 {object} model.ErrorResponse "Erro interno"
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	filter := make(map[string]interface{})

	// Lê query params (status, priority)
	status := r.URL.Query().Get("status")
	if status != "" {
		filter["status"] = status
	}

	priority := r.URL.Query().Get("priority")
	if priority != "" {
		filter["priority"] = priority
	}

	log.WithFields(log.Fields{
		"handler": "ListTasks",
		"filters": filter,
	}).Info("listando tarefas")

	tasks, err := h.service.ListTasks(filter)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "ListTasks",
			"error":   err.Error(),
		}).Error("erro ao listar tarefas")

		http.Error(w, err.Error(), http.StatusInternalServerError) // retorna 500 em caso de erro
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

// ------------------- GET TASK BY ID -------------------

// @Summary Buscar tarefa por ID
// @Description Retorna uma tarefa específica pelo ID
// @Tags tarefas
// @Accept json
// @Produce json
// @Param id path string true "ID da Tarefa"
// @Success 200 {object} model.Task
// @Failure 404 {object} model.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //pega os parametros da url
	id := params["id"]    //extrai o valor do parametro id

	// Validação do ID antes de chamar o service
	if _, err := uuid.Parse(id); err != nil {
		log.WithFields(log.Fields{
			"handler": "GetTask",
			"task_id": id,
		}).Warn("ID inválido recebido")

		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// GET TASK BY ID
	task, err := h.service.GetTask(id) //faz a busca pelo id, a service aplica as regras e chama o repository para pegar os dados do banco
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "GetTask",
			"task_id": id,
			"error":   err.Error(),
		}).Error("erro ao buscar task")

		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, repository.ErrTaskNotFound.Error(), http.StatusNotFound) // 404
			return
		}
		http.Error(w, "erro interno do servidor", http.StatusInternalServerError) // 500
		return
	}

	log.WithFields(log.Fields{
		"handler": "GetTask",
		"task_id": id,
	}).Info("tarefa encontrada")

	json.NewEncoder(w).Encode(task)
}

// ------------------- UPDATE TASK -------------------

// @Summary Atualizar tarefa
// @Description Atualiza o título ou status de uma tarefa pelo ID
// @Tags tarefas
// @Accept json
// @Produce json
// @Param id path string true "ID da Tarefa"
// @Param task body model.UpdateTaskRequest true "Dados da tarefa"
// @Success 200 {object} model.Task
// @Failure 400 {object} model.ErrorResponse "Título inválido (3-100 caracteres)"
// @Failure 404 {object} model.ErrorResponse "Task não encontrada"
// @Failure 400 {object} model.ErrorResponse "Status inválido"
// @Failure 400 {object} model.ErrorResponse "Prioridade inválida"
// @Failure 400 {object} model.ErrorResponse "Data de vencimento no passado"
// @Failure 500 {object} model.ErrorResponse "Erro interno"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Validação do ID antes de chamar o service
	if _, err := uuid.Parse(id); err != nil {
		log.WithFields(log.Fields{
			"handler": "UpdateTask",
			"task_id": id,
		}).Warn("ID inválido recebido")

		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req model.UpdateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithFields(log.Fields{
			"handler": "UpdateTask",
			"task_id": id,
			"error":   err.Error(),
		}).Error("erro ao decodificar JSON")

		http.Error(w, "dados inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask := model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	// Chama o service para atualizar
	if err := h.service.UpdateTask(id, &updatedTask); err != nil {
		log.WithFields(log.Fields{
			"handler": "UpdateTask",
			"task_id": id,
			"error":   err.Error(),
		}).Error("erro ao atualizar tarefa")

		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, repository.ErrTaskNotFound.Error(), http.StatusNotFound) // 404
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest) // 400 para validações
		return
	}

	log.WithFields(log.Fields{
		"handler": "UpdateTask",
		"task_id": id,
	}).Info("tarefa atualizada")

	json.NewEncoder(w).Encode(updatedTask)
}

// ------------------- DELETE TASK -------------------

// @Summary Deletar tarefa
// @Description Deleta uma tarefa específica pelo ID
// @Tags tarefas
// @Accept json
// @Produce json
// @Param id path string true "ID da Tarefa"
// @Success 204
// @Failure 404 {object} model.ErrorResponse "Task não encontrada"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	//Validação do ID antes de chamar o service
	if _, err := uuid.Parse(id); err != nil {
		log.WithFields(log.Fields{
			"handler": "DeleteTask",
			"task_id": id,
		}).Warn("ID inválido recebido")

		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(id); err != nil {
		log.WithFields(log.Fields{
			"handler": "DeleteTask",
			"task_id": id,
			"error":   err.Error(),
		}).Error("erro ao deletar tarefa")

		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"handler": "DeleteTask",
		"task_id": id,
	}).Info("tarefa deletada")

	w.WriteHeader(http.StatusNoContent)
}

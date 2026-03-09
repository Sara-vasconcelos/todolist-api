package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"todolist-api/internal/model"
	"todolist-api/internal/repository"
	"todolist-api/internal/service"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
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
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task //recebe os dados enviados pelo cliente no corpo da requisição

	/*Pega o JSON enviado e converte em uma struct do tipo Task*/
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "dados inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Define status padrão se não tiver sido fornecido
	if task.Status == "" {
		task.Status = "pending"
	}

	/*Chama o TaskService para criar a task, passando um ponteiro para task
	Valida as regras de negocio antes de salvar no banco*/
	if err := h.service.CreateTask(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //retorna 400 em caso de erro
		return
	}

	w.WriteHeader(http.StatusCreated) //retorna 201  em caso de sucesso
	json.NewEncoder(w).Encode(task)   //envia o JSON da task criada de volta para o cliente, incluindo id, created_at e updated_at.
}

// ------------------- LIST TASKS -------------------
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

	tasks, err := h.service.ListTasks(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // retorna 500 em caso de erro
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

// ------------------- GET TASK BY ID -------------------
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) //pega os parametros da url
	id := params["id"]    //extrai o valor do parametro id

	// Validação do ID antes de chamar o service
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// GET TASK BY ID
	task, err := h.service.GetTask(id) //faz a busca pelo id, a service aplica as regras e chama o repository para pegar os dados do banco
	if err != nil {
		log.Printf("erro ao buscar task: %v", err)
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, repository.ErrTaskNotFound.Error(), http.StatusNotFound) // 404
			return
		}
		http.Error(w, "erro interno do servidor", http.StatusInternalServerError) // 500
		return
	}

	json.NewEncoder(w).Encode(task)
}

// ------------------- UPDATE TASK -------------------
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Validação do ID antes de chamar o service
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var updatedTask model.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "dados inválidos: "+err.Error(), http.StatusBadRequest) //400
		return
	}

	// Chama o service para atualizar
	if err := h.service.UpdateTask(id, &updatedTask); err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, repository.ErrTaskNotFound.Error(), http.StatusNotFound) // 404
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest) // 400 para validações
		return
	}

	json.NewEncoder(w).Encode(updatedTask)
}

// ------------------- DELETE TASK -------------------
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	//Validação do ID antes de chamar o service
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(id); err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
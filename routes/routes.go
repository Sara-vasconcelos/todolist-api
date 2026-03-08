

package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"todolist-api/internal/handler"
)

// RegisterRoutes registra todas as rotas da API e retorna o router
func RegisterRoutes(taskHandler *handler.TaskHandler) *mux.Router {
	router := mux.NewRouter()

	// Rotas de tasks
	router.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")          // Criar tarefa
	router.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")             // Listar tarefas
	router.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")          // Buscar tarefa por ID
	router.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")       // Atualizar tarefa
	router.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")    // Deletar tarefa

	// Rota de teste de saúde da API
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is running"))
	}).Methods("GET")

	return router
}
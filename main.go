package main

import (
	"log"
	"net/http"

	"todolist-api/database"
	"todolist-api/internal/handler"
	"todolist-api/internal/logger"
	"todolist-api/internal/repository"
	"todolist-api/internal/service"
	"todolist-api/routes"
)

func main() {
	//Conectar ao MongoDB
	database.ConnectMongo()
	log.Println("Conexão com MongoDB estabelecida")

	//Criar repository, service e handler
	taskRepo := repository.NewTaskRepository()         // acesso ao banco
	taskService := service.NewTaskService(taskRepo)    // regras de negócio
	taskHandler := handler.NewTaskHandler(taskService) // handlers HTTP

	//Registrar rotas
	router := routes.RegisterRoutes(taskHandler)

	logger.Init()

	//Iniciar servidor HTTP
	port := ":8080"
	log.Printf("Servidor rodando na porta %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}

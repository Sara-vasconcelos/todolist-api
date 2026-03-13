package main

// @title Todo List API
// @version 1.0
// @description API para gerenciamento de tarefas
// @host localhost:8080
// @BasePath /

import (
    "log"
    "net/http"

    "todolist-api/database"
    "todolist-api/docs" 
    "todolist-api/internal/handler"
    "todolist-api/internal/logger"
    "todolist-api/internal/repository"
    "todolist-api/internal/service"
    "todolist-api/routes"

    httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
    database.ConnectMongo()

    taskRepo := repository.NewTaskRepository()
    taskService := service.NewTaskService(taskRepo)
    taskHandler := handler.NewTaskHandler(taskService)

    docs.SwaggerInfo.Title = "Todo List API"
    docs.SwaggerInfo.Description = "API para gerenciamento de tarefas"
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = "localhost:8080"
    docs.SwaggerInfo.BasePath = "/"
    docs.SwaggerInfo.Schemes = []string{"http"}

    router := routes.RegisterRoutes(taskHandler)

    //Rota do Swagger com URL 
    router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
        httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
    ))

    logger.Init()

    port := ":8080"
    log.Printf("Servidor rodando na porta %s", port)
    log.Printf("Swagger disponível em: http://localhost%s/swagger/index.html", port)

    if err := http.ListenAndServe(port, router); err != nil {
        log.Fatal("Erro ao iniciar servidor:", err)
    }
}
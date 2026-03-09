package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"todolist-api/database"
	"todolist-api/internal/model"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(m *testing.M) {

	mongoURI := "mongodb://localhost:27017"
	database.ConnectMongo(mongoURI)

	code := m.Run()

	os.Exit(code)
}

func setupTestRepository(t *testing.T) *TaskRepository {

	collection := database.GetCollection("tasks_test_db", "tasks")

	// limpa a collection antes de cada teste
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		t.Fatalf("erro ao limpar banco de teste: %v", err)
	}

	return &TaskRepository{
		collection: collection,
	}
}

func createTestTask() *model.Task {
	return &model.Task{
		Title:       "Task Test",
		Description: "Descrição teste",
		Status:      "pending",
		Priority:    "medium",
		DueDate:     time.Now().Add(24 * time.Hour),
	}
}

func TestCreateTask(t *testing.T) {

	repo := setupTestRepository(t)

	task := createTestTask()

	err := repo.CreateTask(task)

	if err != nil {
		t.Fatalf("erro ao criar task: %v", err)
	}

	if task.ID == "" {
		t.Error("esperava UUID gerado")
	}

	if task.CreatedAt.IsZero() {
		t.Error("esperava CreatedAt preenchido")
	}

	if task.UpdatedAt.IsZero() {
		t.Error("esperava UpdatedAt preenchido")
	}
}

func TestGetTasks(t *testing.T) {

	repo := setupTestRepository(t)

	task := createTestTask()
	_ = repo.CreateTask(task)

	tasks, err := repo.GetTasks(map[string]interface{}{})

	if err != nil {
		t.Fatalf("erro ao buscar tasks: %v", err)
	}

	if len(tasks) == 0 {
		t.Error("esperava ao menos uma task")
	}
}

func TestGetTasksWithFilter(t *testing.T) {

	repo := setupTestRepository(t)

	task := createTestTask()
	task.Status = "completed"

	_ = repo.CreateTask(task)

	filter := map[string]interface{}{
		"status": "completed",
	}

	tasks, err := repo.GetTasks(filter)

	if err != nil {
		t.Fatalf("erro ao buscar tasks com filtro: %v", err)
	}

	if len(tasks) == 0 {
		t.Error("esperava task com status completed")
	}
}

func TestGetTaskByID(t *testing.T) {

	repo := setupTestRepository(t)

	task := createTestTask()
	_ = repo.CreateTask(task)

	result, err := repo.GetTaskByID(task.ID)

	if err != nil {
		t.Fatalf("erro ao buscar task por ID: %v", err)
	}

	if result.Title != task.Title {
		t.Error("task retornada diferente da esperada")
	}
}

func TestGetTaskByIDNotFound(t *testing.T) {

	repo := setupTestRepository(t)

	id := uuid.New().String()

	_, err := repo.GetTaskByID(id)

	if err == nil {
		t.Error("esperava erro de task não encontrada")
	}
}

func TestUpdateTask(t *testing.T) {

	repo := setupTestRepository(t)

	task := createTestTask()
	_ = repo.CreateTask(task)

	update := &model.Task{
		Title:       "Task Atualizada",
		Description: "Descrição Atualizada",
		Status:      "pending",
		Priority:    "high",
		DueDate:     time.Now().Add(48 * time.Hour),
	}

	err := repo.UpdateTask(task.ID, update)

	if err != nil {
		t.Fatalf("erro ao atualizar task: %v", err)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {

	repo := setupTestRepository(t)

	id := uuid.New().String()

	update := createTestTask()

	err := repo.UpdateTask(id, update)

	if err == nil {
		t.Error("esperava erro de task não encontrada")
	}
}

func TestDeleteTask(t *testing.T) {

	repo := setupTestRepository(t)

	task := createTestTask()
	_ = repo.CreateTask(task)

	err := repo.DeleteTask(task.ID)

	if err != nil {
		t.Fatalf("erro ao deletar task: %v", err)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {

	repo := setupTestRepository(t)

	id := uuid.New().String()

	err := repo.DeleteTask(id)

	if err == nil {
		t.Error("esperava erro de task não encontrada")
	}
}

func TestCreateTaskGeneratesValidUUID(t *testing.T) {

	repo := setupTestRepository(t)

	task := createTestTask()

	_ = repo.CreateTask(task)

	_, err := uuid.Parse(task.ID)

	if err != nil {
		t.Error("ID gerado não é um UUID válido")
	}
}
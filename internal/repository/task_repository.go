package repository

import (
	"context"
	"errors"
	"time" //manipula tempo e duração

	"todolist-api/database"
	"todolist-api/internal/model"

	"go.mongodb.org/mongo-driver/bson" //formato de documentos usado pelo MongoDB
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrTaskNotFound = errors.New("task não encontrada")

/*
TaskRepository é responsável pelo acesso ao banco
Estrutura que representa o repositório de tasks
Vai guardar tudo que precisa para acessar o banco de tarefas.
*/
type TaskRepository struct {
	collection *mongo.Collection //Colection:nome do campo dentro da struct , mongo.Collection:Tipo de campo
}

// NewTaskRepository cria uma instância do repository
/*Quem chamar essa função vai receber uma referência para o repositório, que permite acessar e modificar a coleção*/
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{ //retorna o endereço, devolve um ponteiro ao inves de uma cópia
		collection: database.GetCollection("tasks_db", "tasks"), //pega a coleção que definimos no banco de dados
	}
}

// CreateTask insere uma nova task no MongoDB
func (r *TaskRepository) CreateTask(task *model.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //contexto de 5s
	defer cancel()

	// Preenche os campos de data de criação e atualização da task.
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, task) //InsertOne: insere a task no MongoDB.
	return err                                  //retorna um erro se houver
}

// GetTasks retorna todas as tasks ou filtra por status/priority
func (r *TaskRepository) GetTasks(filter map[string]interface{}) ([]model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	/*
		Constrói um mapa de filtros (bson.M) baseado nos parâmetros passados.
		Se filter estiver vazio, retorna todos
	*/
	query := bson.M{}          //bson.M: é o filtro da query, dizendo quais documentos eu quero encontrar/ query := bson.M: inicializa o mapa vazio.
	for k, v := range filter { //percorre o mapa filter com a chave e valor
		query[k] = v //adiciona a chave e valor ao mapa query
	}

	/*Find retorna um cursor para iterar sobre os resultados
	defer cursor.Close() garante que o cursor será fechado
	*/
	cursor, err := r.collection.Find(ctx, query) //busca documentos no Mongo com base no filtro, que definimos acima.
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	//Itera sobre cada documento do cursor, decodifica para Task e adiciona na lista.
	var tasks []model.Task
	for cursor.Next(ctx) {
		var task model.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskByID retorna uma task pelo ID
func (r *TaskRepository) GetTaskByID(id string) (*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Converte string do ID para ObjectID, que é o tipo usado pelo Mongo
	objID, _ := primitive.ObjectIDFromHex(id) // já é válido

	//FindOne busca a task com _id igual ao objID
	//Decode converte o documento Mongo para Task.
	var task model.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

// UpdateTask atualiza uma task pelo ID
func (r *TaskRepository) UpdateTask(id string, updatedTask *model.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id) // já é válido / validação ja foi feita no handler

	// Atualiza a data de modificação.
	updatedTask.UpdatedAt = time.Now()

	/*
		UpdateOne atualiza apenas os campos definidos em $set.
		Retorna erro se não conseguir atualizar
	*/
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"status":      updatedTask.Status,
			"priority":    updatedTask.Priority,
			"due_date":    updatedTask.DueDate,
			"updated_at":  updatedTask.UpdatedAt,
		}},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}

// DeleteTask remove uma task pelo ID
func (r *TaskRepository) DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id) // já é válido
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID}) //DeleteOne remove o documento pelo _id.
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}
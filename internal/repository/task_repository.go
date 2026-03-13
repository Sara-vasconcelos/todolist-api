package repository

import (
	"context"
	"errors"
	"time" //manipula tempo e duração

	"todolist-api/database"
	"todolist-api/internal/model"

	"github.com/google/uuid" // biblioteca para gerar UUID
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"  //formato de documentos usado pelo MongoDB
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

	// Gera um UUID para a task (alteração para não depender do ObjectID do Mongo)
	task.ID = uuid.New().String()

	// Preenche os campos de data de criação e atualização da task.
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "CreateTask",
		"task_id":    task.ID,
	}).Info("inserindo nova task no banco")

	_, err := r.collection.InsertOne(ctx, task) //InsertOne: insere a task no MongoDB.
	if err != nil {
		log.WithFields(log.Fields{
			"repository": "TaskRepository",
			"operation":  "CreateTask",
			"task_id":    task.ID,
			"error":      err.Error(),
		}).Error("erro ao inserir task no MongoDB")

		return err
	}

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "CreateTask",
		"task_id":    task.ID,
	}).Info("task inserida com sucesso")

	return nil //retorna um erro se houver
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

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "GetTasks",
		"filters":    query,
	}).Info("buscando tasks no banco")

	/*Find retorna um cursor para iterar sobre os resultados
	defer cursor.Close() garante que o cursor será fechado
	*/
	cursor, err := r.collection.Find(ctx, query) //busca documentos no Mongo com base no filtro, que definimos acima.
	if err != nil {
		log.WithFields(log.Fields{
			"repository": "TaskRepository",
			"operation":  "GetTasks",
			"error":      err.Error(),
		}).Error("erro ao buscar tasks")

		return nil, err
	}
	defer cursor.Close(ctx)

	//Itera sobre cada documento do cursor, decodifica para Task e adiciona na lista.
	var tasks []model.Task
	for cursor.Next(ctx) {
		var task model.Task
		if err := cursor.Decode(&task); err != nil {
			log.WithFields(log.Fields{
				"repository": "TaskRepository",
				"operation":  "GetTasks",
				"error":      err.Error(),
			}).Error("erro ao decodificar task")

			return nil, err
		}
		tasks = append(tasks, task)
	}

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "GetTasks",
		"total":      len(tasks),
	}).Info("tasks retornadas com sucesso")

	return tasks, nil
}

// GetTaskByID retorna uma task pelo ID
func (r *TaskRepository) GetTaskByID(id string) (*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "GetTaskByID",
		"task_id":    id,
	}).Info("buscando task por id")

	var task model.Task
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&task)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			log.WithFields(log.Fields{
				"repository": "TaskRepository",
				"operation":  "GetTaskByID",
				"task_id":    id,
			}).Warn("task não encontrada")

			return nil, ErrTaskNotFound
		}

		log.WithFields(log.Fields{
			"repository": "TaskRepository",
			"operation":  "GetTaskByID",
			"task_id":    id,
			"error":      err.Error(),
		}).Error("erro ao buscar task")

		return nil, err
	}

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "GetTaskByID",
		"task_id":    id,
	}).Info("task encontrada")

	return &task, nil
}

// UpdateTask atualiza uma task pelo ID
func (r *TaskRepository) UpdateTask(id string, updatedTask *model.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Atualiza a data de modificação.
	updatedTask.UpdatedAt = time.Now()

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "UpdateTask",
		"task_id":    id,
	}).Info("atualizando task")

	/*
		UpdateOne atualiza apenas os campos definidos em $set.
		Retorna erro se não conseguir atualizar
	*/
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": id}, // agora usamos o campo id (UUID)
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

		log.WithFields(log.Fields{
			"repository": "TaskRepository",
			"operation":  "UpdateTask",
			"task_id":    id,
			"error":      err.Error(),
		}).Error("erro ao atualizar task")

		return err
	}

	if result.MatchedCount == 0 {

		log.WithFields(log.Fields{
			"repository": "TaskRepository",
			"operation":  "UpdateTask",
			"task_id":    id,
		}).Warn("task não encontrada para atualização")

		return ErrTaskNotFound
	}

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "UpdateTask",
		"task_id":    id,
	}).Info("task atualizada com sucesso")

	return nil
}

// DeleteTask remove uma task pelo ID
func (r *TaskRepository) DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "DeleteTask",
		"task_id":    id,
	}).Info("removendo task")

	// DeleteOne remove o documento pelo campo id (UUID)
	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {

		log.WithFields(log.Fields{
			"repository": "TaskRepository",
			"operation":  "DeleteTask",
			"task_id":    id,
			"error":      err.Error(),
		}).Error("erro ao deletar task")

		return err
	}

	if result.DeletedCount == 0 {

		log.WithFields(log.Fields{
			"repository": "TaskRepository",
			"operation":  "DeleteTask",
			"task_id":    id,
		}).Warn("task não encontrada para deleção")

		return ErrTaskNotFound
	}

	log.WithFields(log.Fields{
		"repository": "TaskRepository",
		"operation":  "DeleteTask",
		"task_id":    id,
	}).Info("task removida com sucesso")

	return nil
}
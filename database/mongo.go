package database

import (
	"context" //usado para controlar tempo de execução e cancelamento
	"log"     //registra logs de erro ou sucesso
	"time"    //manipula tempo e duração

	"go.mongodb.org/mongo-driver/mongo"         //MongoDB para Go, usado para criar o cliente e interagir com o banco.
	"go.mongodb.org/mongo-driver/mongo/options" //define opções de configuração para o cliente, como URI e timeouts
)

var Client *mongo.Client //armazena a conexão com o banco.

// Essa função cria e inicializa a conexão com o banco.
/*Ela cria um contexto,conecta ao banco,testa a conexão e salva o client global*/
func ConnectMongo(uri string) {

	/*
		context.WithTimeout cria um contexto com limite de tempo para a operação de 10s
		cancel: função que cancela o contexto e libera recursos após o uso
		defer cancel(): garante que a função cancel() será chamada quando ConnectMongo terminar, evitando vazamento de recursos.*/
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	/*.ApplyURI(uri) - preenche todas as informações que constam na URI*/
	clientOptions := options.Client().ApplyURI(uri) //cria uma estrutura de configuração para se conectar.
	/*
	   mongo.Connect cria a conexão com o banco
	   Retorna dois valores:
	   client: objeto de conexão com o Mongo.
	   err: se houve algum erro na conexão.
	   log.Fatal imprime o erro e encerra o programa caso não consiga conectar*/
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}
	/*client.Ping envia um ping para o MongoDB para garantir que ele está acessível.
	  Se der erro o programa encerra*/
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB não respondeu:", err)
	}

	Client = client //Atribui a conexão Mongo à variável global Client. Outras funções podem usar essa mesma conexão sem precisar criar outra.

	log.Println("Conectado ao MongoDB com sucesso")
}

/*
GetCollection retorna uma coleção de tasks do banco todolist
Recebe nome do banco e nome da coleção como parâmetros
Retorna um ponteiro para mongo.Collection, que pode ser usado para fazer inserções, consultas, atualizações, etc..
*/
func GetCollection(databaseName string, collectionName string) *mongo.Collection {
	collection := Client.Database(databaseName).Collection(collectionName)
	return collection
}

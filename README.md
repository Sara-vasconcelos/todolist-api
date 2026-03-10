# ToDo List API

API RESTful para gerenciamento de tarefas (ToDo List), construída em **Golang** com **MongoDB**, seguindo boas práticas de CRUD e regras de negócio definidas.

---

## Tecnologias Utilizadas

- **Back-end:** Golang  
- **Banco de Dados:** MongoDB (NoSQL)  
- **Containerização:** Docker + Docker Compose  
- **Testes/Documentação:** Swagger/OpenAPI

---

## Funcionalidades

- CRUD completo de tarefas: `criar, listar, buscar por ID, atualizar e deletar.` 
- **Filtros:** status (`pending`, `in_progress`, `completed`, `cancelled`) e prioridade (`low`, `medium`, `high`).  

- **Validações de negócio:**
  - Título obrigatório (3 a 100 caracteres).  
  - Status deve ser válido.  
  - Data de vencimento não pode ser no passado.  
  - Tarefas `completed` não podem ser editadas (apenas deletadas).  

---

## Modelo de Dados

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "status": "pending|in_progress|completed|cancelled",
  "priority": "low|medium|high",
  "due_date": "date",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Setup e Execução com Docker

###  1 - Clonar o repositório

```bash
git clone <URL_DO_REPOSITORIO>
cd <NOME_DO_REPOSITORIO>
 
 ```
### 2 - Dar permissão de execução ao script de espera do Mongo

```bash
chmod +x wait-for-mongo.sh
```

O script `wait-for-mongo.sh` garante que o MongoDB esteja pronto antes da API iniciar.

### 3 - Rodar o Docker Compose

- **Linux/Mac:**

```bash
sudo docker compose up --build
```

- **Windows:**

```bash
docker compose up --build
```

### Observação:

`--build`: força a reconstrução das imagens Docker.

`docker compose up`: inicia todos os serviços definidos (api e mongo).

A API ficará disponível em: `http://localhost:8080`

O MongoDB ficará disponível em: `mongodb://localhost:27017`

### 4 - Parar a aplicação

- **Linux/Mac:**

```bash
sudo docker compose down
```
- **Windows:**

```bash
docker compose down
```

## Endpoints Interativos (Testes com curl)

- ### **Criar Tarefa**

```bash

curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{
  "title": "Estudar Golang",
  "description": "Revisar conceitos de goroutines",
  "priority": "high",
  "due_date": "2026-04-10"
}'

```

**Exemplo de retorno:**

```bash

{
  "id": "192f3abe-2178-407b-af7c-25de48415711",
  "title": "Estudar Golang",
  "description": "Revisar conceitos de goroutines",
  "status": "pending",
  "priority": "high",
  "due_date": "2026-04-10",
  "created_at": "2026-03-10T01:33:42Z",
  "updated_at": "2026-03-10T01:33:42Z"
}

```
- ### **Listar Tarefas**

Todas as tarefas: 

```bash
curl -X GET http://localhost:8080/tasks
```

Tarefas com status pending: 

```bash
curl -X GET "http://localhost:8080/tasks?status=pending"
```

Tarefas com prioridade alta: 

```bash
curl -X GET "http://localhost:8080/tasks?priority=high"
```


**Exemplo de retorno:**

```bash
[
  {
    "id": "1a2b3c4d",
    "title": "Estudar Golang",
    "description": "Revisar conceitos de goroutines",
    "status": "pending",
    "priority": "high",
    "due_date": "2026-04-10",
    "created_at": "2026-03-10T01:33:42Z",
    "updated_at": "2026-03-10T01:33:42Z"
  }
]
```

- ### Buscar Tarefa por ID

```bash
curl -X GET http://localhost:8080/tasks/<ID_DA_TAREFA>
```

**Exemplo de retorno:**

```bash
{
  "id": "1a2b3c4d",
  "title": "Estudar Golang",
  "description": "Revisar conceitos de goroutines",
  "status": "pending",
  "priority": "high",
  "due_date": "2026-04-10",
  "created_at": "2026-03-10T01:33:42Z",
  "updated_at": "2026-03-10T01:33:42Z"
}

```

- ### **Atualizar Tarefa**

```bash
curl -X PUT http://localhost:8080/tasks/<ID_DA_TAREFA> \
-H "Content-Type: application/json" \
-d '{
  "title": "Estudar Golang - Atualizado",
  "status": "in_progress"
}'
```

**Exemplo de retorno:**

```bash

{
  "id": "1a2b3c4d",
  "title": "Estudar Golang - Atualizado",
  "description": "Revisar conceitos de goroutines",
  "status": "in_progress",
  "priority": "high",
  "due_date": "2026-04-10",
  "created_at": "2026-03-10T01:33:42Z",
  "updated_at": "2026-03-10T01:50:00Z"
}

```

- ### **Deletar Tarefa**

```bash
curl -X DELETE http://localhost:8080/tasks/<ID_DA_TAREFA>
# Retorno esperado: StatusCode 204 (No Content)
```

## Documentação da API (Swagger)

A documentação interativa está disponível em:

`http://localhost:8080/swagger/index.html_`

**Nela é possível:**

- visualizar todos os endpoints

- testar requisições

- visualizar exemplos de JSON

- verificar códigos de resposta da API

## Observações finais:

- Substitua `<ID_DA_TAREFA>` pelo id retornado ao criar ou listar tarefas.
- Todos os serviços rodam localmente em containers, garantindo que a API funcione da mesma forma em qualquer máquina com Docker.
- Para testar rapidamente, use os exemplos de curl direto no terminal.



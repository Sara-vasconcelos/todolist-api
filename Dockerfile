# imagem oficial do Go
FROM golang:1.22

# pasta de trabalho dentro do container
WORKDIR /app

# Instala netcat (OpenBSD) para o wait-for-mongo.sh
RUN apt-get update && apt-get install -y netcat-openbsd && rm -rf /var/lib/apt/lists/*
# copia os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# baixa as dependências
RUN go mod download

# copia todo o projeto
COPY . .

# compila a aplicação
RUN go build -o /app/main

# expõe a porta da API
EXPOSE 8080

# Usando o script para esperar o Mongo
CMD ["./wait-for-mongo.sh"]
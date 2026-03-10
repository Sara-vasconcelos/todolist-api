#!/bin/sh
# Espera até que o Mongo esteja disponível
# Uso: ./wait-for-mongo.sh host port comando_a_executar

HOST=$1
PORT=$2
shift 2
CMD="$@"

echo "Esperando o MongoDB em $HOST:$PORT..."

# Loop até o Mongo responder
while ! nc -z $HOST $PORT; do
  echo "MongoDB não disponível ainda... esperando 1s"
  sleep 1
done

echo "MongoDB está pronto! Iniciando aplicação..."
exec $CMD
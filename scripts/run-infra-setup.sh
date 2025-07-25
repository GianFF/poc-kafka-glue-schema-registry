#!/bin/bash
set -e

# Solicita el token como argumento o interactivo si no se pasa
if [ -z "$1" ]; then
  read -sp "Ingrese su LOCALSTACK_AUTH_TOKEN: " LOCALSTACK_AUTH_TOKEN
  echo
else
  LOCALSTACK_AUTH_TOKEN="$1"
fi
export LOCALSTACK_AUTH_TOKEN

docker-compose up -d
docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --create --topic users.signedup --partitions 1 --replication-factor 1 || true
docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --list

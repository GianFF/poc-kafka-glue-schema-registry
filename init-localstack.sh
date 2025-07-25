#!/bin/bash
exec > /tmp/localstack-init.log 2>&1
set -e

# Esperar a que LocalStack Glue esté listo
until awslocal --endpoint-url=http://localhost:4566 glue list-registries; do
  echo "Esperando LocalStack Glue..."
  sleep 3
done

echo "Creando Glue Schema Registry y esquema UserSignedUp..."

# Crear registry
awslocal glue create-registry --registry-name user-events

AVRO_SCHEMA=$(cat /schemas/UserSignedUp.avsc | tr -d '\n')

# Crear el esquema en Glue
awslocal glue create-schema \
  --registry-id RegistryName=user-events \
  --schema-name UserSignedUp \
  --data-format AVRO \
  --schema-definition "$AVRO_SCHEMA"

echo "Inicialización completada."
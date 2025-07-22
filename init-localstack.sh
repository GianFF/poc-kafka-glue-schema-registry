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

# Definir el esquema Avro (debe ser una sola línea)
# TOOD: porque estamos definiendo el esquema si ya temeos el asyncapi.yaml?
#   Glue necesita que el esquema esté registrado en su sistema para poder validar, serializar y deserializar mensajes Avro (o JSON Schema) en tiempo real.
#   El registro en Glue se hace usando el formato Avro o JSON Schema, y el registro debe ser creado mediante comandos (CLI, SDK, etc.), no solo documentado.
#   Podrías automatizar la conversión de AsyncAPI → Avro, pero actualmente no existe una integración directa y automática entre AsyncAPI y AWS Glue Schema Registry. Por eso, ambos pasos son necesarios y complementarios.

AVRO_SCHEMA='{
  "type": "record",
  "name": "UserSignedUp",
  "namespace": "com.example",
  "fields": [
    {"name": "user_id", "type": "string"},
    {"name": "email", "type": "string"},
    {"name": "timestamp", "type": "string"}
  ]
}'

# Crear el esquema en Glue
awslocal glue create-schema \
  --registry-id RegistryName=user-events \
  --schema-name UserSignedUp \
  --data-format AVRO \
  --schema-definition "$AVRO_SCHEMA"

echo "Inicialización completada."
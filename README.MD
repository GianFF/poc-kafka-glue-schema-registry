# POC Kafka & AWS Glue Schema Registry

Quiero que generes una prueba de concepto (POC) completa para Kafka con AWS Glue Schema Registry y AsyncAPI que incluya los siguientes componentes:

1. Docker Compose que levante:
- Apache Kafka
- LocalStack con soporte para AWS Glue Schema Registry
- Un archivo AsyncAPI (asyncapi.yaml) que documente un evento de ejemplo llamado UserSignedUp, con un payload simple (user_id, email, timestamp), y que especifique:
  - El topic de Kafka (users.signedup)
  - El productor (auth-service)
  - El consumidor (email-service)
2. Un script init-localstack.sh que:
  - Cree el esquema UserSignedUp que define el AsyncAPI en AWS Glue Schema Registry (usando Avro o JSON Schema)
  - Configure el topic en Kafka (users.signedup) con 1 partición
  - Opcional: registre el esquema en el topic si es posible
3. Una app Node.js (JavaScript) que:
  - Se conecte a Kafka
  - Consuma el topic users.signedup
  - Deserialice el evento usando el esquema del Glue Schema Registry
  - Loguee el contenido del evento en consola
  - Un archivo package.json con las dependencias necesarias para Kafka y AWS Glue Schema Registry (por ejemplo, kafkajs, aws-sdk, avro-js, etc.)

4. Toda la estructura del proyecto organizada en carpetas:
```
poc-kafka-glue-node/
├── docker-compose.yml
├── asyncapi/
│   └── asyncapi.yaml
├── scripts/
│   └── init-localstack.sh
├── app/
│   ├── package.json
│   └── index.js
```


# Recursos
- [LocalStack Schema Evolution](https://docs.localstack.cloud/aws/tutorials/schema-evolution-glue-msk/)
- [AWS Glue Schema Registry](https://docs.aws.amazon.com/glue/latest/dg/schema-registry.html)
- [AsyncAPI](https://www.asyncapi.com/)
- [KafkaJS](https://kafka.js.org/)
- [LocalStack](https://localstack.cloud/)
- [Avro](https://avro.apache.org/)
- [JSON Schema](https://json-schema.org/)


## POC startup 
- npm install
- docker-compose up -d
- docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --create --topic users.signedup --partitions 1 --replication-factor 1
- npm start
- docker-compose down -v

- check topics: `docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --list`

### Consideraciones
1. en el init-localstack.sh es necesario construir el AVRO_SCHEMA, pero no se hace automaticamente con alguna opción del asyncapi o del glue schema registry
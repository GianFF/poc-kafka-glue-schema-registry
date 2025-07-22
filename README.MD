# POC Kafka & AWS Glue Schema Registry

# Recursos
- [LocalStack Schema Evolution](https://docs.localstack.cloud/aws/tutorials/schema-evolution-glue-msk/)
- [AWS Glue Schema Registry](https://docs.aws.amazon.com/glue/latest/dg/schema-registry.html)
- [AsyncAPI](https://www.asyncapi.com/)
- [KafkaJS](https://kafka.js.org/)
- [LocalStack](https://localstack.cloud/)
- [Avro](https://avro.apache.org/)
- [JSON Schema](https://json-schema.org/)

## Startup 
- install dependencies: `npm install`
- start containers: `docker-compose up -d`
- create topic: `docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --create --topic users.signedup --partitions 1 --replication-factor 1`
- start app: `npm start`
- check topics: `docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --list`
- stop containers: `docker-compose down -v`

### Consideraciones
1. en el init-localstack.sh es necesario construir el AVRO_SCHEMA, pero no se hace automaticamente con alguna opción del asyncapi o del glue schema registry
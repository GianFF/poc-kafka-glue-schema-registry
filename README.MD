# POC Kafka & AWS Glue Schema Registry

# Resources
- [LocalStack Schema Evolution](https://docs.localstack.cloud/aws/tutorials/schema-evolution-glue-msk/)
- [AWS Glue Schema Registry](https://docs.aws.amazon.com/glue/latest/dg/schema-registry.html)
- [AsyncAPI](https://www.asyncapi.com/)
- [KafkaJS](https://kafka.js.org/)
- [LocalStack](https://localstack.cloud/)
- [Avro](https://avro.apache.org/)
- [JSON Schema](https://json-schema.org/)

## Running

### Infraestructure
- start containers: `docker-compose up -d`
- create topic: `docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --create --topic users.signedup --partitions 1 --replication-factor 1`
- check topics: `docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --list`
- stop containers: `docker-compose down -v`

### Javascript
- `cd javascript`
- install dependencies: `npm install`
- run: `npm run start`

### Java
- `cd java`
- install dependencies: `./mvnw clean install`
- build: `./mvnw clean package`
- run: `java -jar target/kafka-glue-poc-1.0.0.jar`

### GoLang
- `cd golang`
- install dependencies: `go mod tidy`
- run: `go run ./cmd/main.go`

## 🐛 Debugging

### Ver logs de LocalStack
```bash
docker-compose logs localstack
```

### Ver mensajes en Kafka
```bash
docker-compose exec kafka kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic users.signedup --from-beginning
```

### Verificar esquemas en Glue
```bash
aws --endpoint-url=http://localhost:4566 glue get-schema-version --schema-id RegistryName=user-events,SchemaName=UserSignedUp --schema-version-number LatestVersion=true
```

### Consideraciones
1. en el init-localstack.sh es necesario construir el AVRO_SCHEMA, pero no se hace automaticamente con alguna opción del asyncapi o del glue schema registry
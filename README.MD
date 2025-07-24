# POC Kafka & AWS Glue Schema Registry

## 📚 Resources
- [LocalStack Schema Evolution](https://docs.localstack.cloud/aws/tutorials/schema-evolution-glue-msk/)
- [AWS Glue Schema Registry](https://docs.aws.amazon.com/glue/latest/dg/schema-registry.html)
- [AsyncAPI](https://www.asyncapi.com/)
- [KafkaJS](https://kafka.js.org/)
- [LocalStack](https://localstack.cloud/)
- [Avro](https://avro.apache.org/)
- [JSON Schema](https://json-schema.org/)

## 🏃 Running

### Infraestructure
```bash
sh scripts/run-infra.sh
```

### Javascript
```bash
sh scripts/run-javascript.sh
```

### Java
```bash
sh scripts/run-java.sh
```

### GoLang
```bash
sh scripts/run-golang.sh
```

## 🐛 Debugging

### Ver logs de LocalStack
```bash
docker-compose logs localstack
```

### Ver mensajes en Kafka
```bash
docker-compose exec kafka kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic users.signedup --from-beginning
```

### Ver offset de Kafka
```bash
docker-compose exec kafka kafka-run-class.sh kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic users.signedup --time -1
```
### Verificar esquemas en Glue
```bash
aws --endpoint-url=http://localhost:4566 glue get-schema-version --schema-id RegistryName=user-events,SchemaName=UserSignedUp --schema-version-number LatestVersion=true
```

## 🧠 Consideraciones

### Definición del esquema Avro

El archivo `init-localstack.sh` contiene la definición del esquema Avro para el evento `UserSignedUp`, el cual está documentado en el archivo `asyncapi.yaml`. Este esquema es necesario para registrar el evento en AWS Glue Schema Registry y garantizar la validación de los mensajes producidos.
Este es el primer problema que podemos observar: duplicamos la definición del esquema.

```bash
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
```

### JavaScript – Mapeo manual al esquema
En `javascript/app/producer.js`, el objeto `user` se serializa usando `avroType.toBuffer(user)`. Esta línea lanza una excepción si el objeto no cumple con el esquema. En una aplicación productiva, el objeto `user` podría tener una estructura diferente a la esperada por el evento, por lo que sería necesario mapearlo explícitamente al esquema del evento.  
Para evitar acoplamientos innecesarios, lo ideal sería autogenerar este mapper a partir de la definición de AsyncAPI.

### GoLang – Mapeo manual al esquema
En `golang/internal/producer/producer.go`, el evento `UserSignedUp` se construye manualmente como un `map[string]interface{}`. Este enfoque también requiere que el desarrollador conozca y respete la estructura del esquema, lo cual puede volverse frágil si el esquema cambia.

### Java – Generación automática del tipo
En `java/src/main/avro/UserSignedUp.avsc` se encuentra la definición del esquema Avro. Gracias a un plugin, se genera automáticamente la clase `UserSignedUp`, lo que permite trabajar con tipos seguros y evitar errores de mapeo manual.  
Sin embargo, aún es necesario mantener sincronizado el esquema Avro con la definición de AsyncAPI.  

## Conclusiones

Tanto para GoLang como para JS fueron necesarias 3 definiciones para el mismo esquema:
- AsyncAPI
- Esquema Avro
- DTO

En cambio en Java solo tenemos 2 definiciones:
- AsyncAPI
- Esquema Avro  

#### GoLang:
```go
	userMap := map[string]interface{}{
		"user_id":   user.UserID,
		"email":     user.Email,
		"timestamp": user.Timestamp,
	}
```     

#### Javascript:
```js
    event = avroType.toBuffer(user);
```     

#### Java:
```java
public void sendUserSignedUpEvent(UserSignedUp user) {
    try {
        Schema schema = glueSchemaClient.getSchema("user-events", "UserSignedUp");
        byte[] payload = AvroUtil.serialize(user, schema);
        kafkaTemplate.send(topic, payload).get();
        ...
```



### Recomendación general
En una arquitectura bien diseñada, los eventos como `UserSignedUp` deberían definirse una sola vez (idealmente en AsyncAPI) y desde allí:

- Autogenerar los esquemas Avro
- Autogenerar los DTOs en los distintos lenguajes
- Autogenerar los mappers si el modelo de dominio difiere del modelo del evento

Esto reduce el acoplamiento, mejora la mantenibilidad y asegura la consistencia entre servicios.
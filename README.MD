# POC Kafka & AWS Glue Schema Registry

Esta prueba de concepto se centra en revisar el AWS Glue Schema Registry y su integración con Kafka.    
La solucion se implemento en 3 lenguajes diferentes (JavaScript, GoLang, Java).   
En cada app, se implementa un productor y un consumidor. El productor fallara el 50% de las veces que quiera enviar un mensaje ya que agregara un typo que no respete el formato del evento segun lo definido en AsyncAPI y AWS Glue Schema Registry.

Se utiliza localstack para simular el entorno de AWS y docker compose para gestionar los servicios de Kafka y AWS Glue Schema Registry.   
Un script init-localstack.sh se encarga de inicializar localstack y registrar el esquema en AWS Glue Schema Registry haciendo uso del AsyncAPI.

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
docker-compose exec localstack awslocal glue list-schemas
```


## 🧠 Consideraciones

### Registro del esquema en AWS Glue Schema Registry

Para registrar el evento en AWS Glue Schema Registry y garantizar la validación de los mensajes producidos necesitamos la definición del esquema en AVRO y enviarlo al registry. 
Si documentamos el esquema en AsyncAPI, y no referenciamos a un archivo externo con la defincion del evento, tendremos que duplicar la definición del esquema.    
He aqui el primer problema que podemos observar.

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
El archivo `init-localstack.sh` toma la definición del esquema en AVRO y lo registra en AWS Glue Schema Registry.
El archivo `schemas/UserSignedUp.avsc` contiene la definición del esquema en AVRO.
El archivo `asyncapi.yaml` contiene la definicion del evento, server y canales. Para la definicion del evento se apalanca del archivo anterior (`schemas/UserSignedUp.avsc`) y de esta forma podemos tener una unica fuente de verdad.


### JavaScript – Mapeo manual al esquema
En `javascript/app/producer.js`, el objeto `user` se serializa usando `avroType.toBuffer(user)`. Esta línea lanza una excepción si el objeto no cumple con el esquema. En una aplicación productiva, el objeto `user` podría tener una estructura diferente a la esperada por el evento, por lo que sería necesario mapearlo explícitamente al esquema del evento.   
Para evitar acoplamientos innecesarios, lo ideal sería autogenerar un mapper y un DTO a partir de la definición de AsyncAPI.

### GoLang – Mapeo manual al esquema
En `golang/internal/producer/producer.go`, el evento `UserSignedUp` se construye manualmente como un `map[string]interface{}`. Este enfoque también requiere que el desarrollador conozca y respete la estructura del esquema, lo cual puede volverse frágil si el esquema cambia.

### Java – Generación automática del tipo
Gracias a un plugin, se genera automáticamente la clase `UserSignedUp`, lo que permite trabajar con tipos seguros y evitar errores de mapeo manual.  

## Conclusiones

Gracias a la definicion del AVRO y su reutilizacion en el AsyncAPI hemos evitado triplicar la fuente de la verdad. En el caso de Java, gracias al plugin, hemos evitado tener que generar manualmente la clase `UserSignedUp`. En el caso de GoLang y JS, debemos generar manualmente el mapper y el DTO.   

En caso que no usaramos la definicion del AVRO en el AsyncAPI, tendriamos este escenario:

Tanto para GoLang como para JS serian necesarias 3 definiciones para el mismo esquema:
- AsyncAPI
- Esquema Avro
- DTO

En cambio en Java solo tendriamos 2 definiciones:
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
Los eventos deberían definirse una sola vez (idealmente en AsyncAPI) y desde allí:

- Autogenerar los esquemas Avro
- Autogenerar los DTOs en los distintos lenguajes
- Autogenerar los mappers si el modelo de dominio difiere del modelo del evento

Esto reduce el acoplamiento, mejora la mantenibilidad y asegura la consistencia entre servicios.

### TODO
- Revisar la documentación de AsyncAPI para ver si hay una forma de autogenerar los esquemas Avro a partir de esa definicion y evitar la separacion en dos archivos. Sobre todo teniendo en cuenta que no queremos tener problemas a la hora de subir los asyncAPI al hub de eventos
- Probar el versionado de eventos en AWS Glue Schema Registry
- Review costos de AWS Glue Schema Registry
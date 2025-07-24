package com.example.kafkagluepoc.kafka;

import com.example.kafkagluepoc.UserSignedUp;
import com.example.kafkagluepoc.schema.GlueSchemaClient;
import com.example.kafkagluepoc.util.AvroUtil;
import org.apache.avro.Schema;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Service
public class KafkaProducerService {
    private final GlueSchemaClient glueSchemaClient;
    private final KafkaTemplate<String, byte[]> kafkaTemplate;

    @Value("${kafka.topic}")
    private String topic;

    public KafkaProducerService(GlueSchemaClient glueSchemaClient, KafkaTemplate<String, byte[]> kafkaTemplate) {
        this.glueSchemaClient = glueSchemaClient;
        this.kafkaTemplate = kafkaTemplate;
    }

    public void sendUserSignedUpEvent(UserSignedUp user) {
        try {
            // 1. Obtener y validar el esquema desde Glue
            Schema schema = glueSchemaClient.getSchema("user-events", "UserSignedUp");

            // 2. Serializar el objeto Avro específico (validación automática)
            byte[] payload = AvroUtil.serialize(user, schema);

            // 3. Enviar a Kafka
            kafkaTemplate.send(topic, payload).get();
            System.out.println("🚀 Mensaje enviado:");
            System.out.println("   Usuario: " + user.getUserId() + ", Email: " + user.getEmail());
        } catch (Exception e) {
            System.err.println("💥 Error enviando mensaje a Kafka: " + e.getMessage());
        }
    }
}

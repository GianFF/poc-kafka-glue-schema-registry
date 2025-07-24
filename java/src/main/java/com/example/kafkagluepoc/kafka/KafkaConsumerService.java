package com.example.kafkagluepoc.kafka;

import com.example.kafkagluepoc.schema.GlueSchemaClient;
import com.example.kafkagluepoc.util.AvroUtil;
import org.apache.avro.Schema;
import org.apache.avro.generic.GenericRecord;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

@Service
public class KafkaConsumerService {
    private final GlueSchemaClient glueSchemaClient;

    public KafkaConsumerService(GlueSchemaClient glueSchemaClient) {
        this.glueSchemaClient = glueSchemaClient;
    }

    @KafkaListener(topics = "#{'${kafka.topic}'}", groupId = "glue-poc-group")
    public void listen(byte[] message) {
        try {
            Schema schema = glueSchemaClient.getSchema("user-events", "UserSignedUp");
            GenericRecord record = AvroUtil.deserialize(message, schema);
            System.out.println("📨 Mensaje recibido: " + record);
        } catch (Exception e) {
            System.err.println("💥 Error al deserializar mensaje: " + e.getMessage());
        }
    }
}

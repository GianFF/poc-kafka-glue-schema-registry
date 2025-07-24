package com.example.kafkagluepoc.schema;

import org.apache.avro.Schema;
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials;
import software.amazon.awssdk.auth.credentials.StaticCredentialsProvider;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.glue.GlueClient;
import software.amazon.awssdk.services.glue.model.GetSchemaVersionRequest;
import software.amazon.awssdk.services.glue.model.GetSchemaVersionResponse;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.util.concurrent.ConcurrentHashMap;

@Component
public class GlueSchemaClient {
    private final ConcurrentHashMap<String, Schema> cache = new ConcurrentHashMap<>();
    private final GlueClient glueClient;

    public GlueSchemaClient(
        @Value("${aws.glue.endpoint}") String endpoint,
        @Value("${aws.glue.region}") String region,
        @Value("${aws.glue.access-key}") String accessKey,
        @Value("${aws.glue.secret-key}") String secretKey
    ) {
        this.glueClient = GlueClient.builder()
            .endpointOverride(java.net.URI.create(endpoint))
            .region(Region.of(region))
            .credentialsProvider(StaticCredentialsProvider.create(
                AwsBasicCredentials.create(accessKey, secretKey)))
            .build();
    }

    public Schema getSchema(String registryName, String schemaName) {
        String cacheKey = registryName + ":" + schemaName;
        if (cache.containsKey(cacheKey)) {
            return cache.get(cacheKey);
        }
        // Aquí deberías obtener la última versión del schema desde Glue
        // Por simplicidad, se asume que el schemaId es schemaName
        GetSchemaVersionRequest req = GetSchemaVersionRequest.builder()
                .schemaId(b -> b.registryName(registryName).schemaName(schemaName))
                .schemaVersionNumber(b -> b.latestVersion(true))
                .build();
        GetSchemaVersionResponse resp = glueClient.getSchemaVersion(req);
        String schemaDefinition = resp.schemaDefinition();
        Schema schema = new Schema.Parser().parse(schemaDefinition);
        cache.put(cacheKey, schema);
        return schema;
    }
}

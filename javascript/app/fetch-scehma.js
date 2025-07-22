const glueClient = require('./aws-glue-client');
const avro = require('avro-js');

const glue = glueClient.init();

// Cache del esquema para evitar consultas repetidas
let cachedAvroType = null;

// Obtiene el esquema desde Glue con cache
async function fetchSchema() {
  if (!cachedAvroType) {
    console.log('📥 Obteniendo esquema desde Glue Schema Registry...');
    
    const schema = await glue.getSchemaVersion({
      SchemaId: { SchemaName: 'UserSignedUp', RegistryName: 'user-events' },
      SchemaVersionNumber: { LatestVersion: true }
    }).promise();

    cachedAvroType = avro.parse(JSON.parse(schema.SchemaDefinition));
    console.log('✅ Esquema cargado y cacheado');
  }
  
  return cachedAvroType;
}

module.exports = { fetchSchema };
const AWS = require('aws-sdk');
const kafkaClient = require('./kafka-client');
const avro = require('avro-js');


const kafka = kafkaClient.init('auth-service');
const topic = 'users.signedup';
const producer = kafka.producer();

function randomUser() {
  const id = Math.floor(Math.random() * 100000);
  return {
    user_id: `user${id}`,
    email: `user${id}@example.com`,
    timestamp: new Date().toISOString()
  };
}

// Configura AWS SDK para LocalStack
const glue = new AWS.Glue({
  endpoint: 'http://localhost:4566', // LocalStack Edge
  region: 'us-east-1',
  accessKeyId: 'test', // Default para LocalStack
  secretAccessKey: 'test'
});

// Obtiene y parsea el esquema desde Glue
async function fetchAvroType() {
  const schema = await glue.getSchemaVersion({
    SchemaId: { SchemaName: 'UserSignedUp', RegistryName: 'user-events' },
    SchemaVersionNumber: { LatestVersion: true }
  }).promise();

  return avro.parse(JSON.parse(schema.SchemaDefinition));
}

const run = async () => {
  const avroType = await fetchAvroType();
  await producer.connect();
  console.log(`Enviando mensajes al topic "${topic}" cada 1 segundo...`);

  setInterval(async () => {
    const event = randomUser();

    let avroBuffer;
    try {
      avroBuffer = avroType.toBuffer(event); // Validación aquí
    } catch (err) {
      console.error('Error de validación Avro:', err);
      return;
    }

    try {
      await producer.send({
        topic,
        messages: [
          { value: avroBuffer }
        ]
      });
      console.log('Mensaje enviado:', event);
    } catch (err) {
      console.error('Error enviando mensaje:', err);
    }
  }, 1000);
};

module.exports = { run };
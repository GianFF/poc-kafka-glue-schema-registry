const kafkaClient = require('./kafka-client');
const userSignedUpSchema = require('./user-schema');

const kafka = kafkaClient.init('email-service');
const topic = 'users.signedup';
const consumer = kafka.consumer({ groupId: 'email-service-group' });

const run = async () => {
  await consumer.connect();
  await consumer.subscribe({ topic, fromBeginning: true });

  console.log(`Esperando mensajes en el topic "${topic}"...`);

  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      try {
        const event = userSignedUpSchema.fromBuffer(message.value);
        console.log('Evento recibido:', event);
      } catch (err) {
        console.error('Error deserializando mensaje:', err);
      }
    }
  });
};

module.exports = { run };
const kafkaClient = require('./kafka-client');
const { fetchSchema } = require('./fetch-scehma');

const kafka = kafkaClient.init('email-service');
const topic = 'users.signedup';
const consumer = kafka.consumer({ groupId: 'email-service-group' });

const run = async () => {
  await consumer.connect();
  await consumer.subscribe({ topic, fromBeginning: true });

  console.log(`🔍 Consumidor esperando mensajes en el topic "${topic}"...`);

  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      const messageInfo = `Topic: ${topic}, Partition: ${partition}, Offset: ${message.offset}`;
      
      try {
        // 1. Obtener esquema desde Glue (con cache)
        const avroType = await fetchSchema();
        
        // 2. Validar y deserializar el mensaje
        const event = avroType.fromBuffer(message.value);
        
        console.log('✅ Mensaje recibido:');
        console.log(`   ${messageInfo}`);
        console.log('   Datos:', event);
        
        // 3. Procesar el evento (aquí iría tu lógica de negocio)
        await processUserSignedUpEvent(event);
        
      } catch (err) {
        console.log('❌ Error procesando mensaje:');
        console.log(`   ${messageInfo}`);
        
        if (err.message.includes('invalid')) {
          console.log('   🚫 Mensaje no cumple con el esquema:', err.message);
        } else {
          console.log('   💥 Error técnico:', err.message);
        }
      }
    }
  });
};

// Simula el procesamiento del evento
async function processUserSignedUpEvent(event) {
  console.log(`   📧 Enviando email de bienvenida a: ${event.email}`);
  console.log(`   👤 Usuario: ${event.user_id} registrado en: ${event.timestamp}`);
}

module.exports = { run };
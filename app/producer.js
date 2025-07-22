const kafkaClient = require('./kafka-client');
const { fetchSchema } = require('./fetch-scehma');


const kafka = kafkaClient.init('auth-service');
const topic = 'users.signedup';
const producer = kafka.producer();

const run = async () => {
  await producer.connect();
  console.log(`🚀 Productor iniciado - enviando mensajes al topic "${topic}" cada 2 segundos...`);

  setInterval(async () => {
    const user = randomUser();
    const [event, error] = await buildEvent(user);

    if (error) {
      console.log(`❌ Mensaje rechazado - no cumple esquema:`);
      console.log(`   Usuario: ${JSON.stringify(user)}`);
      console.log(`   Error: ${error.message}`);
      return;
    }

    await publish(event, user);
  }, 2000);
};

async function publish(event, user) {
  try {
    await producer.send({
      topic,
      messages: [
        { value: event }
      ]
    });
    console.log('✅ Mensaje enviado:');
    console.log(`   Usuario: ${user.user_id}, Email: ${user.email}`);
  } catch (err) {
    console.log('💥 Error enviando mensaje a Kafka:', err.message);
  }
}

async function buildEvent(user) {
  let event;
  try {
    const avroType = await fetchSchema();
    // VALIDACIÓN DEL ESQUEMA: aquí es donde se valida el mensaje antes del envío
    event = avroType.toBuffer(user); // Si el objeto no cumple el esquema, lanza excepción
    return [event, null];
  } catch (err) {
    // Error de validación - el mensaje no cumple con el esquema
    return [null, err];
  } 
}

function randomUser() {
  const id = Math.floor(Math.random() * 100000);

  let email = `user${id}@example.com`;
  const user_id = `user${id}`;
  const timestamp = new Date().toISOString();

  if (id % 2 === 0){
    // Intencionalmente inválido
    email = undefined;
  }

  return { user_id, email, timestamp };
}

module.exports = { run };
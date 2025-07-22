const avro = require('avro-js');

// TODO: La integración directa con el Glue Schema Registry en Node.js no es tan transparente como en Java, pero para la POC (y usando LocalStack), puedes usar el esquema Avro localmente para deserializar los mensajes.
// Definir el esquema Avro (debe ser igual al registrado en Glue)
const userSignedUpSchema = avro.parse({
  type: 'record',
  name: 'UserSignedUp',
  namespace: 'com.example',
  fields: [
    { name: 'user_id', type: 'string' },
    { name: 'email', type: 'string' },
    { name: 'timestamp', type: 'string' }
  ]
});

module.exports = userSignedUpSchema;

const { Kafka } = require('kafkajs');

const init = (clientId) => new Kafka({
  clientId,
  brokers: ['localhost:9092'] 
});

module.exports = { init };

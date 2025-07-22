const { Kafka } = require('kafkajs');

const init = (clientId) => new Kafka({
  clientId,
  brokers: ['localhost:9093'] 
});

module.exports = { init };

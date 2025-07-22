const consumer = require('./app/consumer');
const producer = require('./app/producer-valid');
const invalidProducer = require('./app/producer-invalid');

Promise.all([ 
  consumer.run(),
  producer.run(),
  invalidProducer.run()
]).catch(console.error);
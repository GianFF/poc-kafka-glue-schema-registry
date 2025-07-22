const consumer = require('./app/consumer');
const producer = require('./app/producer');

Promise.all([ 
  consumer.run(),
  producer.run(),
]).catch(console.error);
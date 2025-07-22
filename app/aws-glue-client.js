const AWS = require('aws-sdk');

const init = () => new AWS.Glue({
  endpoint: 'http://localhost:4566', // LocalStack Edge
  region: 'us-east-1',
  accessKeyId: 'test', // Default para LocalStack
  secretAccessKey: 'test'
});

module.exports = { init };

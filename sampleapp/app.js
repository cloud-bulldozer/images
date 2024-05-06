const http = require('http');
const os = require('os');


var handler = function(request, response) {
  console.log("Received request from " + request.connection.remoteAddress);
  response.writeHead(200);
};

var www = http.createServer(handler);
www.listen(8080);

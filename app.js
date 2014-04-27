// Generated by CoffeeScript 1.7.1
(function() {
  var fs, http, log, mongo, mongoClient, mongoUri, port, server;

  http = require('http');

  fs = require('fs');

  mongo = require('mongodb');

  log = console.log;

  mongoUri = process.env.MONGOHQ_URL || 'mongodb://localhost/mydb';

  mongoClient = mongo.MongoClient;

  mongoClient.connect(mongoUri, function(err, db) {
    var collection;
    if (err) {
      throw error;
    }
    return collection = db.collection('test');
  });

  server = http.createServer(function(req, res) {
    var data, output;
    data = {
      "Content-type": "text/html"
    };
    res.writeHead(200, data);
    output = fs.readFileSync("html/index.html", "utf-8");
    return res.end(output);
  });

  port = process.env.PORT || 5000;

  server.listen(port, function() {
    return console.log("http server listening on port " + server.address().port);
  });

}).call(this);

var http = require ('http')
//var socketio = require('socket.io');
var fs = require ('fs');
var redis = require('redis');

var server = http.createServer(function (req, res) {
	res.writeHead(200, {"Content-type": "text/html"});
	var output = fs.readFileSync("html/index.html", "utf-8");
	res.end(output);
})

var port = process.env.PORT || 5000; // Use the port that Heroku provides or default to 5000
server.listen(port, function() {
  console.log("http server listening on port %d ", server.address().port);
});

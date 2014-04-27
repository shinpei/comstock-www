http = require 'http'
fs = require 'fs'
mongo = require 'mongodb'
log = console.log

mongoUri = process.env.MONGOHQ_URL || 'mongodb://localhost/mydb'
mongoClient = mongo.MongoClient;

mongoClient.connect(mongoUri, (err, db) ->
    throw error if err
    collection = db.collection 'test'
)

server = http.createServer (req, res)->
    data  =
        "Content-type" : "text/html"
        
    res.writeHead(200, data);
    output = fs.readFileSync("html/index.html", "utf-8");
    res.end(output);

port = process.env.PORT || 5000;
server.listen(port, ->
    console.log "http server listening on port " + server.address().port;
)    

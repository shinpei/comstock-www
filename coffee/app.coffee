http = require 'http'
fs = require 'fs'
mongo = require 'mongodb'
log = console.log
url = require 'url'
S = require 'string'
querystring = require 'querystring'
uuid = require 'node-uuid'
socketio = require 'socket.io'


mongoUri = process.env.MONGOHQ_URL || 'mongodb://localhost/mydb'
mongoClient = mongo.MongoClient;

DOCROOT = "documents"
COLLECTION_NAME = "test"

getHandler = (filepath, req, res) ->
    fs.readFile(filepath, "utf-8", (err, data) ->
        throw err if err
        header =
            "Content-type" : ""
        if S(filepath).endsWith(".css")
            header["Content-type"] = "text/css";
        else if S(filepath).endsWith(".html")
            header["Content-type"] = "text/html"
        else if S(filepath).endsWith(".js")
            header["Content-type"] = "application/javascript"
        res.writeHead(200, header);
        res.end(data);
    );

postCommand = (command, user, date, desc) ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(COLLECTION_NAME)
        log 'removing documents'
        id = uuid.v1()
        collection.remove ((err, result) ->
            throw err if err
            log "colelction cleared!"
            oneData =
                "id" : id
                "date": date
                "data":
                    "command": command
                    "user" : user
                    "desc" : desc
                
            collection.insert(oneData, (err, docs) ->
                throw err if err
                log "Just inserted, " + docs.length
                collection.find({}).toArray (err, docs) ->
                    throw err if err
                    docs.forEach (doc) ->
                        log "found document:" + doc.data.command
            )
        )
    )


getCommand = () ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(COLLECTION_NAME)
        docs = collection.find({}, limit: 5)
        docs.count (err, count) ->
            log " #{count} document(s) fund"
            log "=========================="
            docs.toArray (err, docs) ->
                throw err if err
                for doc in docs then log doc
                    
    )

        
server = http.createServer (req, res) ->
    filepath = ''
    isIgnore = false;
    pathname = url.parse(req.url).pathname;
    log "pathname=" + pathname
    if pathname == '/'
        filepath = DOCROOT + "/index.html"
        getHandler(filepath, req, res);
    else if pathname == '/favicon.ico'
        res.writeHead(404);
        return;
    else if pathname == "/postCommand"
        query = url.parse(req.url).query
        params = querystring.parse(query);
        postCommand(params.command, params.user, params.date, params.desc);
        res.writeHead(200)
    else if pathname == "/getCommand"
        res.writeHead(200, {"Content-type": "plain/text"})
    else if pathname == "/search"
        res.writeHead(200, {"Content-type": "plain/text"})
    else
        filepath = DOCROOT + req.url;
        getHandler(filepath, req, res);
        return;

## io = socketio.listen(server);
port = process.env.PORT || 5000;
server.listen(port, ->
    log "http server listening on port " + server.address().port;
)

###
io.configure ->
    io.set("transports", ["xhr-polling"]);
    io.set("polling duration", 10);


io.sockets.on('connection', (socket) ->
    socket.on('fetchCommands', (data) ->
        commandData = getCommnad();
        io.sockets.emit('recvCommand', {data: commandData});
    )
)
###

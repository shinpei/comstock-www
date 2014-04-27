http = require 'http'
fs = require 'fs'
mongo = require 'mongodb'
log = console.log
url = require 'url'
S = require 'string'
querystring = require 'querystring'

mongoUri = process.env.MONGOHQ_URL || 'mongodb://localhost/mydb'
mongoClient = mongo.MongoClient;

DOCROOT = "documents"

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
        collection = db.collection('test')
        console.log 'removing documents'
        collection.remove ((err, result) ->
            throw err if err
            console.log "colelction cleared!"
            oneData =
                "id" :
                    "command": command
                    "user" : user
                    "date" : date
                    "desc" : desc
                
            collection.insert(oneData, (err, docs) ->
                throw err if err
                console.log "Just inserted, " + docs.length
                collection.find({}).toArray (err, docs) ->
                    throw err if err
                    docs.forEach (doc) ->
                        console.log "found document:" + doc.id.command
            )
        )
    )


server = http.createServer (req, res)->
    filepath = ''
    isIgnore = false;
    pathname = url.parse(req.url).pathname;
    console.log "pathname=" + pathname
    if pathname == '/'
        filepath = DOCROOT + "/index.html"
        getHandler(filepath, req, res);
        return;
    else if pathname == '/favicon.ico'
        res.writeHead(404);
        return;
    else if pathname == "/postCommand"
        query = url.parse(req.url).query
        params = querystring.parse(query);
        console.log params
        postCommand(params.command, params.user, params.date, params.desc);
        res.writeHead(200)
    else
        filepath = DOCROOT + req.url;
        getHandler(filepath, req, res);
        return;


port = process.env.PORT || 5000;
server.listen(port, ->
    console.log "http server listening on port " + server.address().port;
)

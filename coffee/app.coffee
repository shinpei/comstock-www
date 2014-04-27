http = require 'http'
fs = require 'fs'
mongo = require 'mongodb'
log = console.log
url = require 'url'

S = require 'string'
querystring = require 'querystring'

mongoUri = process.env.MONGOHQ_URL || 'mongodb://localhost/mydb'
mongoClient = mongo.MongoClient;

postHandler = (req, res, postData) ->
    console.log querystring.parse(postData).text
    console.log req


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


            
mongoClient.connect(mongoUri, (err, db) ->
    throw error if err
    collection = db.collection 'test'
)
DOCROOT = "documents"

server = http.createServer (req, res)->
    filepath = ''
    isIgnore = false;
    postData = '';
    
    if req.url == '/'
        filepath = DOCROOT + "/index.html"
    else if req.url == '/favicon.ico'
        isIgnore = true;
    else
        filepath = DOCROOT + req.url;

    console.log "Request: " + filepath;
    req.addEventListner('data', (chunk) ->
        postData = chunk
    )
    
    if isIgnore == true
        res.writeHead(404)
        return
    
    if req.method == "GET"
        getHandler(filepath, req, res)
    else if req.method == "POST"
        postHandler(req, res, postData)
    

port = process.env.PORT || 5000;
server.listen(port, ->
    console.log "http server listening on port " + server.address().port;
)

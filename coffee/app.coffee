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
USER_COLLECTION = "user"
AUTH_COLLECTION = "authinfo" # "auth" is reserve words for mongo client
DATA_COLLECTION = "test"

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
        else if S(filepath).endsWith(".ico")
            header["Content-type"] = "image/x-icon"
        else if S(filepath).endsWith("png")
            header["Conent-type"] = "image/png"
        res.writeHead(200, header);
        res.end(data);
    );

postCommand = (command, user, date, desc) ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(DATA_COLLECTION)
#        log 'removing documents'
        id = uuid.v1()
#        collection.remove ((err, result) ->
#            throw err if err
#            log "colelction cleared!"
        oneData =
            "id" : id
            "user": user
            "date": date
            "data":
                "command": command
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


makeHTMLResponse = (msg, status) ->
    response = '<html><head><!-- Loading Bootstrap --><link href="css/bootstrap.min.css" rel="stylesheet"><!-- Loading Flat UI --><link href="css/flat-ui.css" rel="stylesheet"><link href="css/demo.css" rel="stylesheet"></head><body>'
    response += msg
    response += "</body></html>"


loginAs = (username, password, res) ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(USER_COLLECTION)
        createdNewUser = false
        uid = 0;
        log "god username?"
        log username
        docs = collection.findOne({mail: username}, (err, item) ->
            throw err if err
            log "User finding?"
            log item
            if item == null
                # cannot find user. register it
                collection.find().count((err, count) ->
                    throw err if err
                    uid = count + 1
                    newUser =
                       "uid": uid
                       "username": ""
                       "mail": username # first register, email is uname
                       "created": 1244
                       "lastLogin":0
                    collection.insert(newUser, (err, docs) ->
                        throw err if err
                        log "uid is " + uid
                        addAuthenticate(uid, password)
                        response = makeHTMLResponse("User added, thank you for registering", 200)
                        res.writeHead(200, {"Content-type": "text/html"});
                        res.end(response)
                    )
                )
            else
                # found user
                log item
                log "User found, now authenticate"
                uid = item.uid;
                authenticate(uid, password, res)
        ) # findOne done
    )
    return

authenticate = (uid, password, res) ->
    log "Authentication process got uid="+ uid
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(AUTH_COLLECTION)
        docs = collection.findOne({uid: uid}, (err, item) ->
            if item != null
                # found uid
                if password == item.password
                    response = makeHTMLResponse("Success")
                    res.writeHead(200,{"Content-type": "text/html"});
                    res.end(response)
                    log "authenticate done with ok"
                else
                    response = makeHTMLResponse("Failed")
                    res.writeHead(403 ,{"Content-type": "text/html"});
                    res.end(response)
                    log "authenticate done with ng"
            else
                response = makeHTMLResponse("Not found")
                res.writeHead(404 ,{"Content-type": "text/html"});
                res.end(response)
                log "cannot find uid"
                
        )
    )
    return 


addAuthenticate = (uid, password) ->
    # add user-password to the db
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(AUTH_COLLECTION)
        oneData =
            "uid": uid
            "date": ""
            "password":password
        collection.insert(oneData, (err, docs) ->
            throw err if err
        )
    )
    log "password insertion done"
    return



writeAsHtml = (doc) ->
    log "Logging.."
    log doc
    output = ""
    output += "<div class='commandContain'>"
    output += '<pre class="prettyprint">' + doc.data.command + "</pre>";
    output += "<span class='desc'>" + doc.data.desc + "</span>";
    output += "<span class='user'> by " + doc.data.user + "</span>";
    output += "</div>"
    return output;

getCommand = (res) ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(COLLECTION_NAME)
        response = '<html><head><!-- Loading Bootstrap --><link href="css/bootstrap.min.css" rel="stylesheet"><!-- Loading Flat UI --><link href="css/flat-ui.css" rel="stylesheet"><link href="css/demo.css" rel="stylesheet"></head><body>'
        docs = collection.find({}, limit: 5)
        docs.each (err, doc) ->
            throw err if err
            if doc == null
                res.writeHead(200, {"Content-type": "text/html"});
                response += "</body></html>"
                log response
                res.end(response);
                return;
            response += writeAsHtml(doc);
        return
    )

        
server = http.createServer (req, res) ->
    filepath = ''
    isIgnore = false;
    pathname = url.parse(req.url).pathname;
    log "pathname=" + pathname
    if pathname == '/'
        filepath = DOCROOT + "/index.html"
        getHandler(filepath, req, res);
    else if pathname == "/postCommand"
        query = url.parse(req.url).query
        params = querystring.parse(query);
        postCommand(params.command, params.user, params.date, params.desc);
        filepath = DOCROOT + "/index.html"
        getHandler(filepath, req, res);
    else if pathname == "/getCommand"
        getCommand(res)
    else if pathname == "/loginAs"
        query = url.parse(req.url).query
        params = querystring.parse(query)
        username = params.mail
        password = params.password

        loginAs(username, password, res);
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

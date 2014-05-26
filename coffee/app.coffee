http = require 'http'
fs = require 'fs'
mongo = require 'mongodb'
log = console.log
url = require 'url'
S = require 'string'
querystring = require 'querystring'
uuid = require 'node-uuid'
socketio = require 'socket.io'


mongoUri = process.env.MONGOHQ_URL || 'mongodb://localhost/comstock-www'
mongoClient = mongo.MongoClient;

DOCROOT = "documents"
USER_COLLECTION = "user"
AUTH_COLLECTION = "authinfo" # "auth" is reserve words for mongo client
DATA_COLLECTION = "commands"
SESSION_COLLECTION = "session"

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

postCommand = (token, command, date, res) ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(SESSION_COLLECTION)
        doc = collection.findOne({token: token}, (err, item) ->
            throw err if err
            if item == null
                # not found
                log "token not found, means, hasn't login"
                response = "Hasn't login yet"
                res.writeHead(404, {"Content-type": "text/html"})
                res.end(response)
            else
                uid = item.uid
                collection = db.collection(DATA_COLLECTION)
                id = uuid.v1()
                cmd = new Command()
                cmd.id =  id
                cmd.uid = uid
                cmd.date = date;
                cmd.data = 
                        "command": command
                        "desc" : ""
                       
                collection.insert(cmd, (err, docs) ->
                    throw err if err
                    log "Just inserted, " + docs.length
                    res.writeHead(200, {"Content-type": "text/html"});
                    res.end()
                )
        )
)

makeHTMLResponse = (msg, status) ->
    response = '<html><head><!-- Loading Bootstrap --><link href="css/bootstrap.min.css" rel="stylesheet"><!-- Loading Flat UI --><link href="css/flat-ui.css" rel="stylesheet"><link href="css/demo.css" rel="stylesheet"></head><body>'
    response += msg
    response += "</body></html>"

loginAs = (user, password, res) ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(USER_COLLECTION)
        createdNewUser = false
        uid = 0;
        docs = collection.findOne({mail: user.mail}, (err, item) ->
            throw err if err
            log "Searching User"
            log item
            if item == null
                # user not found
                response = makeHTMLResponse("Not Found")
                res.writeHead(404, {"Content-type": "text/html"});
                res.end(response)
            else
                # found user
                log item
                log "User found, now authenticate"
                uid = item.uid;
                # check weather it's already logged in
                collection = db.collection(SESSION_COLLECTION)
                collection.findOne({uid: uid}, (err, item) ->
                    throw err if err
                    if item == null
                        #Couldn't find user, proceed authenticate
                        authenticate(uid, password, res)
                    else
                        #Already logged in, return "already loggedin"
                        response = makeHTMLResponse("Conflict")
                        res.writeHead(409, {"Content-type": "text/html"});
                        res.end(response)
                )
        ) # findOne done
    )
    return

loginOrRegister = (user, password, res) ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        collection = db.collection(USER_COLLECTION)
        createdNewUser = false
        uid = 0;
        docs = collection.findOne({mail: user.mail}, (err, item) ->
            throw err if err
            log "User finding?"
            log item
            if item == null
                # cannot find user. register it
                collection.find().count((err, count) ->
                    throw err if err
                    date = new Date();
                    user.uid = count + 1
                    user.created = date.getTime()
                    user.lastLogin = date.getTime();
                    collection.insert(user, (err, docs) ->
                        throw err if err
                        log "uid is " + user.uid
                        addAuthenticate(user.uid, password)
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
                    accessToken = uuid.v1()
                    registerToken(db, uid, accessToken)
                    response = accessToken
                    res.writeHead(200,{"Content-type": "text/html"});
                    res.end(response)
                    log "authenticate success!"
                else
                    response = makeHTMLResponse("Login Denied")
                    res.writeHead(403 ,{"Content-type": "text/html"});
                    res.end(response)
                    log "authentication denied for wrong password"
            else
                response = makeHTMLResponse("Not found")
                res.writeHead(404 ,{"Content-type": "text/html"});
                res.end(response)
                log "Authentication defnied because user uid not found"
        )
    )


registerToken = (db, uid, token) ->
    collection = db.collection(SESSION_COLLECTION)
    ses = new Session(token, uid)
    log "Registering session" + ses.token

    collection.insert(ses, (err, docs) ->
        throw err if err
    )


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

listCommands = (token, res) ->
    mongoClient.connect(mongoUri, (err, db) ->
        throw err if err
        log "mongo connected"
        collection = db.collection(SESSION_COLLECTION)
        doc = collection.findOne({token: token}, (err, item) ->
            throw err if err
            if item == null
                # not found
                log "token not found, means, hasn't login"
                response = "Hasn't login yet"
                res.writeHead(404, {"Content-type": "text/html"});
                res.end(response)
            else
                # found session, continue
                log "session found!"
                collection = db.collection(DATA_COLLECTION)
                docs = collection.find({uid: item.uid }, limit: 100)
                response = ""
                responseObjs = []
                docs.each (err, doc) ->
                    throw err if err
                    if doc == null
                        res.writeHead(200, {"Content-type": "text/html"});
                        response = JSON.stringify responseObjs
                        log response
                        res.end(response);
                        return;
                    docObj =
                        Cmd : doc.data.command
                        Timestamp: doc.date
                    responseObjs.push(docObj)

            return
        )
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
        token = params.authinfo
        postCommand(token, params.cmd, params.date, res);
    else if pathname == "/list"
        query = url.parse(req.url).query
        params = querystring.parse(query)
        token = params.authinfo
        listCommands(token, res)
    else if pathname == "/loginOrRegister"
        query = url.parse(req.url).query
        params = querystring.parse(query)
        mail = params.mail
        password = params.password
        user = new User(mail, "", "")
        loginOrRegister(user, password, res);
    else if pathname == "/loginAs"
        query = url.parse(req.url).query
        params = querystring.parse(query)
        mail = params.mail
        password = params.password
        user = new User(mail, "", "")
        loginAs(user, password, res);
    else if pathname == "/search"
        res.writeHead(200, {"Content-type": "plain/text"})
    else
        filepath = DOCROOT + req.url;
        getHandler(filepath, req, res);
        return;

port = process.env.PORT || 5000;
server.listen(port, ->
    log "http server listening on port " + server.address().port;
)

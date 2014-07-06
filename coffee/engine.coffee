class Engine
    constructor: ->
        return

    registerUser: (user, password, res) ->
        mongoClient.connect(mongoUri, (err, db) ->
            throw err if err
            collection = db.collection(USER_COLLECTION)
            createdNewUser = false
            uid = 0;
            docs = collection.findOne({mail: user.mail}, (err, item) ->
                throw err if err
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
                            collection = db.collection(AUTH_COLLECTION)
                            oneData =
                                "uid": user.uid
                                "date": ""
                                "password":password
                             collection.insert(oneData, (err, docs) ->
                                 throw err if err
                                 db.close()
                                 response =
                                     message: "User added, thank you for registering"
                                 res.writeHead(200, {"Content-type": "application/json"});
                                 res.end(JSON.stringify(response))
                            )
                        )
                    )
                else
                    log "User found, you cannot create duplicated user"
                    db.close()
                    response =
                        message: "It's already registered email. Please try another one, or if you don't know about it, please let us know"
                    res.writeHead(401, {"Content-type": "application/json"})
                    res.end(JSON.stringify(response))
            ) # findOne done
        )
    checkSession: (token, res) ->
        mongoClient.connect(mongoUri, (err, db) ->
            throw err if err
            collection = db.collection(SESSION_COLLECTION)
            doc = collection.findOne({token: token}, (err, item) ->
                throw err if err
                if item == null
                    # not found
                    log "token not found, means, hasn't login"
                    db.close()
                    response =
                        message: "Hasn't login yet"
                    res.writeHead(404, {"Content-type": "application/json"})
                    res.end(JSON.stringify(response))
                else
                    dateobj = new Date();
                    if item.expires < dateobj.getTime()
                        # session expires
                        response =
                            message: "Session expires, please login again"
                        res.writeHead(500, {"Content-type": "application/json"})
                        res.end(JSON.stringify(response))
                        cleanupSession(db, collection, token);
                    else
                        response =
                            message:"Session is alive"
                        res.writeHead(200, {"Content-type": "application/json"})
                        res.end()
                        db.close()
            )
        )
                            
    cleanupSession: (db, collection, token) ->
        collection.remove({token: token}, (err, item) ->
            throw err if err
            db.close()
        )

            
    deleteUser: (user, res) ->
        mongoClient.connect(mongoUri, (err, db) ->
            throw err if err
            collection = db.collection(USER_COLLECTION)
            doc = collection.findOne({mail:user.mail}, (err, item) ->
                throw err if err
                if item == null
                    db.close()
                    response = "User not found"
                    res.writeHead(404, {"Content-type": "text/html"});
                    res.end(response)
                else
                    uid = parseInt item.uid;
                    collection = db.collection(DATA_COLLECTION)
                    collection.remove({uid: uid}, (err, num) ->
                        throw err if err
                        db.close()
                        response = "delete done"
                        res.writeHead(200, {"Content-type": "text/html"})
                        res.end(response)
                    )
            )
        )


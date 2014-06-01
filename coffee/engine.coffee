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
                                response = makeHTMLResponse("User added, thank you for registering", 200)
                                res.writeHead(200, {"Content-type": "text/html"});
                                res.end(response)
                                db.close()
                            )
                        )
                    )
                else
                    log "User found, you cannot create duplicated user"
                    response = makeHTMLResponse("It's already registered email. Please try another one, or if you don't know about it, please let us know")
                    res.writeHead(401, {"Content-type": "text/html"})
                    res.end(response)
                    db.close()
            ) # findOne done
        )

    addAuthenticate : (db, res, uid, password) ->
                
    deleteUser: (user, res) ->
        mongoClient.connect(mongoUri, (err, db) ->
            throw err if err
            collection = db.collection(USER_COLLECTION)
            doc = collection.findOne({mail:user.mail}, (err, item) ->
                throw err if err
                if item == null
                    response = "User not found"
                    res.writeHead(404, {"Content-type": "text/html"});
                    res.end(response)
                    db.close()
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


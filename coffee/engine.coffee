class Engine
    constructor: ->
        return
        
    deleteUser: (user, res) ->
        log "hi"
        mongoClient.connect(mongoUri, (err, db) ->
            throw err if err
            collection = db.collection(USER_COLLECTION)
            doc = collection.findOne({mail:user.mail}, (err, item) ->
                throw err if err
                if item == null
                    log "user not found"
                    response = "User not found"
                    res.writeHead(404, {"Content-type": "text/html"});
                    res.end(response)
                else
                    log "user found: ", item.uid
                    uid = parseInt item.uid;
                    collection = db.collection(DATA_COLLECTION)
                    collection.remove({uid: uid}, (err, num) ->
                        log "removinfg :", num
                        throw err if err
                        db.close()
                        response = "delete done"
                        res.writeHead(200, {"Content-type": "text/html"})
                        res.end(response)
                    )
            )
        )


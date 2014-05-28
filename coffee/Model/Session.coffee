class Session
    "token": ""
    "uid" : 0
    "expires" : 0
        
    constructor: (token, uid) ->
        @token = token
        @uid = uid
        date = new Date();
        #exp = date.getTime() + 10800 # 3hours
        exp = date.getTime() + 1 # 3hours
        @expires = exp
        return

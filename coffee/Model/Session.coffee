class Session
    "token": ""
    "uid" : 0
    "expires" : 0
    SECOND: 1000
    MINUTE: 60 * @SECOND
    HOUR : 60 * @MINUTE        
    constructor: (token, uid) ->
        @token = token
        @uid = uid
        date = new Date();
        exp = date.getTime() + 24 * @HOUR # 24hours
        @expires = exp
        return

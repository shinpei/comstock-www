// Generated by CoffeeScript 1.7.1
(function() {
  var AUTH_COLLECTION, Command, DATA_COLLECTION, DOCROOT, S, SESSION_COLLECTION, Session, USER_COLLECTION, User, addAuthenticate, authenticate, fs, getHandler, http, listCommands, log, loginAs, loginOrRegister, makeHTMLResponse, mineTypes, mongo, mongoClient, mongoUri, path, port, postCommand, querystring, registerToken, server, url, uuid, writeAsHtml;

  http = require('http');

  fs = require('fs');

  mongo = require('mongodb');

  log = console.log;

  url = require('url');

  S = require('string');

  querystring = require('querystring');

  uuid = require('node-uuid');

  path = require('path');

  mongoUri = process.env.MONGOHQ_URL || 'mongodb://localhost/comstock-www';

  mongoClient = mongo.MongoClient;

  DOCROOT = "documents";

  USER_COLLECTION = "user";

  AUTH_COLLECTION = "authinfo";

  DATA_COLLECTION = "commands";

  SESSION_COLLECTION = "session";

  mineTypes = {
    ".html": "text/html",
    ".css": "text/css",
    ".js": "application/javascript",
    ".png": "image/png",
    ".gif": "image/gif",
    ".ico": "image/x-icon",
    ".jpg": "image/jpeg"
  };

  getHandler = function(filepath, req, res) {
    return fs.exists(filepath, function(exists) {
      if (exists) {
        return fs.readFile(filepath, function(err, data) {
          var ext, header;
          if (err) {
            throw err;
          }
          ext = path.extname(filepath);
          header = {
            "Content-type": mineTypes[ext],
            "Content-length": data.length
          };
          res.writeHead(200, header);
          return res.end(data);
        });
      }
    });
  };

  postCommand = function(token, command, date, res) {
    return mongoClient.connect(mongoUri, function(err, db) {
      var collection, doc;
      if (err) {
        throw err;
      }
      collection = db.collection(SESSION_COLLECTION);
      return doc = collection.findOne({
        token: token
      }, function(err, item) {
        var cmd, id, response, uid;
        if (err) {
          throw err;
        }
        if (item === null) {
          log("token not found, means, hasn't login");
          response = "Hasn't login yet";
          res.writeHead(404, {
            "Content-type": "text/html"
          });
          return res.end(response);
        } else {
          uid = item.uid;
          collection = db.collection(DATA_COLLECTION);
          id = uuid.v1();
          cmd = new Command();
          cmd.id = id;
          cmd.uid = uid;
          cmd.date = date;
          cmd.data = {
            "command": command,
            "desc": ""
          };
          return collection.insert(cmd, function(err, docs) {
            if (err) {
              throw err;
            }
            log("Just inserted, " + docs.length);
            res.writeHead(200, {
              "Content-type": "text/html"
            });
            return res.end();
          });
        }
      });
    });
  };

  makeHTMLResponse = function(msg, status) {
    var response;
    response = '<html><head><!-- Loading Bootstrap --><link href="css/bootstrap.min.css" rel="stylesheet"><!-- Loading Flat UI --><link href="css/flat-ui.css" rel="stylesheet"><link href="css/demo.css" rel="stylesheet"></head><body>';
    response += msg;
    return response += "</body></html>";
  };

  loginAs = function(user, password, res) {
    mongoClient.connect(mongoUri, function(err, db) {
      var collection, createdNewUser, docs, uid;
      if (err) {
        throw err;
      }
      collection = db.collection(USER_COLLECTION);
      createdNewUser = false;
      uid = 0;
      return docs = collection.findOne({
        mail: user.mail
      }, function(err, item) {
        var response;
        if (err) {
          throw err;
        }
        log("Searching User");
        log(item);
        if (item === null) {
          response = makeHTMLResponse("Not Found");
          res.writeHead(404, {
            "Content-type": "text/html"
          });
          return res.end(response);
        } else {
          log(item);
          log("User found, now authenticate");
          uid = item.uid;
          collection = db.collection(SESSION_COLLECTION);
          return collection.findOne({
            uid: uid
          }, function(err, item) {
            if (err) {
              throw err;
            }
            if (item === null) {
              return authenticate(uid, password, res);
            } else {
              response = item.token;
              res.writeHead(409, {
                "Content-type": "text/html"
              });
              return res.end(response);
            }
          });
        }
      });
    });
  };

  loginOrRegister = function(user, password, res) {
    mongoClient.connect(mongoUri, function(err, db) {
      var collection, createdNewUser, docs, uid;
      if (err) {
        throw err;
      }
      collection = db.collection(USER_COLLECTION);
      createdNewUser = false;
      uid = 0;
      return docs = collection.findOne({
        mail: user.mail
      }, function(err, item) {
        if (err) {
          throw err;
        }
        log("User finding?");
        log(item);
        if (item === null) {
          return collection.find().count(function(err, count) {
            var date;
            if (err) {
              throw err;
            }
            date = new Date();
            user.uid = count + 1;
            user.created = date.getTime();
            user.lastLogin = date.getTime();
            return collection.insert(user, function(err, docs) {
              var response;
              if (err) {
                throw err;
              }
              log("uid is " + user.uid);
              addAuthenticate(user.uid, password);
              response = makeHTMLResponse("User added, thank you for registering", 200);
              res.writeHead(200, {
                "Content-type": "text/html"
              });
              return res.end(response);
            });
          });
        } else {
          log(item);
          log("User found, now authenticate");
          uid = item.uid;
          return authenticate(uid, password, res);
        }
      });
    });
  };

  authenticate = function(uid, password, res) {
    log("Authentication process got uid=" + uid);
    return mongoClient.connect(mongoUri, function(err, db) {
      var collection, docs;
      if (err) {
        throw err;
      }
      collection = db.collection(AUTH_COLLECTION);
      return docs = collection.findOne({
        uid: uid
      }, function(err, item) {
        var accessToken, response;
        if (item !== null) {
          if (password === item.password) {
            accessToken = uuid.v1();
            registerToken(db, uid, accessToken);
            response = accessToken;
            res.writeHead(200, {
              "Content-type": "text/html"
            });
            res.end(response);
            return log("authenticate success!");
          } else {
            response = makeHTMLResponse("Login Denied");
            res.writeHead(403, {
              "Content-type": "text/html"
            });
            res.end(response);
            return log("authentication denied for wrong password");
          }
        } else {
          response = makeHTMLResponse("Not found");
          res.writeHead(404, {
            "Content-type": "text/html"
          });
          res.end(response);
          return log("Authentication defnied because user uid not found");
        }
      });
    });
  };

  registerToken = function(db, uid, token) {
    var collection, ses;
    collection = db.collection(SESSION_COLLECTION);
    ses = new Session(token, uid);
    log("Registering session" + ses.token);
    return collection.insert(ses, function(err, docs) {
      if (err) {
        throw err;
      }
    });
  };

  addAuthenticate = function(uid, password) {
    mongoClient.connect(mongoUri, function(err, db) {
      var collection, oneData;
      if (err) {
        throw err;
      }
      collection = db.collection(AUTH_COLLECTION);
      oneData = {
        "uid": uid,
        "date": "",
        "password": password
      };
      return collection.insert(oneData, function(err, docs) {
        if (err) {
          throw err;
        }
      });
    });
    log("password insertion done");
  };

  writeAsHtml = function(doc) {
    var output;
    log("Logging..");
    output = "";
    output += "<div class='commandContain'>";
    output += '<pre class="prettyprint">' + doc.data.command + "</pre>";
    output += "<span class='desc'>" + doc.data.desc + "</span>";
    output += "<span class='user'> by " + doc.data.user + "</span>";
    output += "</div>";
    return output;
  };

  listCommands = function(token, res) {
    return mongoClient.connect(mongoUri, function(err, db) {
      var collection, doc;
      if (err) {
        throw err;
      }
      log("mongo connected");
      collection = db.collection(SESSION_COLLECTION);
      doc = collection.findOne({
        token: token
      }, function(err, item) {
        var docs, response, responseObjs;
        if (err) {
          throw err;
        }
        if (item === null) {
          log("token not found, means, hasn't login");
          response = "Hasn't login yet";
          res.writeHead(404, {
            "Content-type": "text/html"
          });
          res.end(response);
        } else {
          log("session found!");
          collection = db.collection(DATA_COLLECTION);
          docs = collection.find({
            uid: item.uid
          }, {
            limit: 100
          });
          response = "";
          responseObjs = [];
          docs.each(function(err, doc) {
            var docObj;
            if (err) {
              throw err;
            }
            if (doc === null) {
              res.writeHead(200, {
                "Content-type": "text/html"
              });
              response = JSON.stringify(responseObjs);
              log(response);
              res.end(response);
              return;
            }
            docObj = {
              Cmd: doc.data.command,
              Timestamp: doc.date
            };
            return responseObjs.push(docObj);
          });
        }
      });
    });
  };

  server = http.createServer(function(req, res) {
    var basename, dirname, isIgnore, mail, params, password, pathname, query, token, user;
    isIgnore = false;
    basename = path.basename(req.url) || 'index.html';
    dirname = path.dirname(req.url);
    if (dirname === "/") {
      dirname = "";
    }
    log(basename);
    if (basename.indexOf("postCommand") === 0) {
      query = url.parse(req.url).query;
      params = querystring.parse(query);
      token = params.authinfo;
      return postCommand(token, params.cmd, params.date, res);
    } else if (basename.indexOf("list") === 0) {
      query = url.parse(req.url).query;
      params = querystring.parse(query);
      token = params.authinfo;
      return listCommands(token, res);
    } else if (basename.indexOf("loginOrRegister") === 0) {
      query = url.parse(req.url).query;
      params = querystring.parse(query);
      mail = params.mail;
      password = params.password;
      user = new User(mail, "", "");
      return loginOrRegister(user, password, res);
    } else if (basename.indexOf("loginAs") === 0) {
      query = url.parse(req.url).query;
      params = querystring.parse(query);
      mail = params.mail;
      password = params.password;
      user = new User(mail, "", "");
      return loginAs(user, password, res);
    } else if (basename.indexOf("search") === 0) {
      return res.writeHead(200, {
        "Content-type": "plain/text"
      });
    } else {
      pathname = dirname + "/" + basename;
      pathname = DOCROOT + pathname;
      log("pathname: " + pathname);
      getHandler(pathname, req, res);
    }
  });

  port = process.env.PORT || 5000;

  server.listen(port, function() {
    return log("http server listening on port " + server.address().port);
  });

  Command = (function() {
    function Command() {}

    Command.prototype["id"] = "";

    Command.prototype["uid"] = "";

    Command.prototype["date"] = "";

    Command.prototype["data"] = {
      "command": "",
      "desc": ""
    };

    return Command;

  })();

  Session = (function() {
    Session.prototype["token"] = "";

    Session.prototype["uid"] = 0;

    function Session(token, uid) {
      this.token = token;
      this.uid = uid;
    }

    return Session;

  })();

  User = (function() {
    User.prototype["mail"] = "";

    User.prototype["username"] = "";

    User.prototype["uid"] = "";

    User.prototype["created"] = "";

    User.prototype["lastLogin"] = "";

    function User(mail, username, uid) {
      this.mail = mail;
      this.username = username;
      this.uid = uid;
    }

    return User;

  })();

}).call(this);

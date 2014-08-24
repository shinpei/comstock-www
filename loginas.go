package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	AUTH_COLLECTION    string = "authinfo"
	USER_COLLECTION    string = "user"
	SESSION_COLLECTON  string = "session"
	COMMAND_COLLECTION string = "commands"
)

// maybe refact to GetUserSession
func GetUserSession(db *mgo.Database, token string) (session *model.Session, err error) {
	c := db.C(SESSION_COLLECTON)
	session = &model.Session{}
	err = c.Find(bson.M{"token": token}).One(&session)
	if err != nil {
		// session not found. reject.
		err = cmodel.ErrSessionNotFound
	}
	// TODO: compare time. document's time is unix time
	//	unixTime := time.Unix(session.Expires, 0)
	//	println(unixTime.Format(time.RFC3339))

	return
}

func LoginAsHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()
	// make sure param exists
	params, _ := url.ParseQuery(req.URL.RawQuery)
	if params["mail"] == nil || params["password"] == nil {
		// error
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("login request mail:%#v, %#v\n", params["mail"], params["mail"][0])
	s, err := loginAs(db, model.CreateLoginRequest(params["mail"][0], params["password"][0]))
	if err == cmodel.ErrUserNotFound || err == cmodel.ErrIncorrectPassword {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if err == cmodel.ErrServerSystem {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	if err == cmodel.ErrAlreadyLogin {
		w.WriteHeader(http.StatusConflict)

	}

	w.Write([]byte(s.Token))
}

func loginAs(db *mgo.Database, l *model.LoginRequest) (s *model.Session, err error) {
	c := db.C(USER_COLLECTION)
	user := model.User{}
	log.Println("l.Mail:", l.Mail())
	err = c.Find(bson.M{"mail": l.Mail()}).One(&user)
	if err != nil {
		log.Println("Counln't find user, ", l.Mail())
		err = cmodel.ErrUserNotFound
		return
	}
	// close db connection??

	c = db.C(SESSION_COLLECTON)
	s = new(model.Session)
	err = c.Find(bson.M{"uid": user.UID}).One(&s)
	if err != nil {
		// session not found. authenticate
		log.Println("Error occured. check it: ", err.Error())
		s, err = authenticateUser(db, user.UID, l)
		if err != nil {
			// if error occurs, s is nil
			s = nil
		}
		// insert session
		err = c.Insert(s)
		if err != nil {
			s = nil
		}

	} else {
		log.Println("Session fonud, for user ", l.Mail())
		// check weather it expires or not
		exp := time.Unix(s.Expires, 0)
		now := time.Now()
		if exp.Before(now) {
			// Session expired!
			// INFO: Made new variable for avoiding shadowing compile error
			newSession, errAuth := authenticateUser(db, user.UID, l)
			// update session
			if errAuth != nil {
				err = errAuth
				return
			}
			errAuth = c.Update(bson.M{"uid": user.UID}, newSession)
			if errAuth != nil {
				s = nil
				err = cmodel.ErrServerSystem
				return
			}
			log.Println("Session found, and expired, but updated")
			s = newSession
		} else {
			log.Println("Session found, but not expired")
		}
		err = cmodel.ErrAlreadyLogin
	}
	return
}

func authenticateUser(db *mgo.Database, uid int, l *model.LoginRequest) (s *model.Session, err error) {
	c := db.C(AUTH_COLLECTION)
	auth := model.Auth{}
	err = c.Find(bson.M{"uid": uid}).One(&auth)
	if err != nil {
		// error occured.
		log.Println("Error occured:", err.Error())
		log.Println("User seems not registered")

	} else {
		// check password
		if auth.Password != l.Pass() {
			err = cmodel.ErrIncorrectPassword
			log.Println("Incorrect password")
			return
		}
		// password ok
		c = db.C(SESSION_COLLECTON)
		s = model.CreateSession(uuid.New(), uid)
	}
	return
}

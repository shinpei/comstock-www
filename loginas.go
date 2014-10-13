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

func GetUserSession(db *mgo.Database, token string) (session *model.Session, err error) {
	c := db.C(SESSION_COLLECTON)
	log.Println("token: ", token)
	session = &model.Session{}
	err = c.Find(bson.M{"token": token}).One(&session)
	if err != nil {
		// session not found. reject.
		err = &cmodel.SessionNotFoundError{}
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
		log.Println("Either mail or password is empty for login request")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("login request mail:%#v, %#v\n", params["mail"], params["mail"][0])
	s, err := loginAs(db, model.CreateLoginRequest(params["mail"][0], params["password"][0]))
	if _, ok := err.(*cmodel.UserNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.IncorrectPasswordError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.ServerSystemError); ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	if _, ok := err.(*cmodel.AlreadyLoginError); ok {
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
		err = &cmodel.UserNotFoundError{}
		return
	}

	c = db.C(SESSION_COLLECTON)
	s = new(model.Session)
	err = c.Find(bson.M{"uid": user.UID}).One(&s)
	if err != nil {
		// session not found. authenticate
		log.Println("Error occured. check it: ", err.Error())
		s, err = authenticateUser(db, user.UID, l, nil)
		if err != nil {
			// if error occurs, s is nil
			s = nil
		}
		// insert session
		err = c.Insert(s)
		if err != nil {
			log.Println("Cannot write session, ", err.Error())
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
			newSession, errAuth := authenticateUser(db, user.UID, l, &s.ID)
			// update session
			if errAuth != nil {
				err = errAuth
				return
			}

			errAuth = c.Update(bson.M{"uid": user.UID}, newSession)
			if errAuth != nil {
				s = nil
				log.Println("Update failed, ", errAuth.Error())
				err = &cmodel.ServerSystemError{}
				return
			}
			log.Println("Session found, and expired, but updated")
			s = newSession
		} else {
			log.Println("Session found, and it's still avlie")
		}
		err = &cmodel.AlreadyLoginError{}
	}
	return
}

func authenticateUser(db *mgo.Database, uid int, l *model.LoginRequest, updateForExistingID *bson.ObjectId) (s *model.Session, err error) {
	c := db.C(AUTH_COLLECTION)
	auth := model.Auth{}
	err = c.Find(bson.M{"uid": uid}).One(&auth)
	if err != nil {
		// error occured.
		log.Println("User seems not registered:", err.Error())
	} else {
		// check password
		if auth.Password != l.Pass() {
			err = &cmodel.IncorrectPasswordError{}
			log.Println("Incorrect password")
			return
		}
		if updateForExistingID != nil {
			s = model.UpdateSessionToken(updateForExistingID, uuid.New(), uid)
		} else {
			s = model.CreateNewSession(uuid.New(), uid)
		}
	}

	return
}

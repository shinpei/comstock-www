package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
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
	// compare time. document's time is unix time
	//	unixTime := time.Unix(session.Expires, 0)
	//	println(unixTime.Format(time.RFC3339))
	println(session.Expires)
	println(time.Now().Format(time.RFC3339))
	println(time.Now().Unix())
	return
}

func LoginAs(db *mgo.Database, l *model.LoginRequest) (s *model.Session, err error) {
	c := db.C(USER_COLLECTION)
	user := model.User{}
	log.Println("SHINPEI!")
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
			return
		}
	} else {
		// TODO: session found. update lastlogin
		log.Println("User already logined")
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
		fmt.Println("Error occured:", err.Error())
		fmt.Println("User seems not registered")
	} else {
		// check password
		if auth.Password != l.Pass() {
			err = cmodel.ErrIncorrectPassword
			return
		}
		// password ok
		println("uuid=", uuid.New())
		c = db.C(SESSION_COLLECTON)
		s = model.CreateSession(uuid.New(), uid)
		c.Insert(s)
	}
	//
	return
}

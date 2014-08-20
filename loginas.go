package main

import (
	"fmt"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
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
	return
}

func LoginAs(db *mgo.Database, l *model.LoginRequest) (s *model.Session, err error) {
	c := db.C(USER_COLLECTION)
	user := model.User{}

	log.Println("l.Mail:", l.Mail())
	err = c.Find(bson.M{"mail": l.Mail()}).One(&user)
	if err != nil {
		log.Println("Counln't find user")
		err = cmodel.ErrUserNotFound
		return
	}
	// found user
	log.Printf("user;%#v\n", user)
	log.Println("mail:'", user.Mail, "'")
	// close db connection??

	c = db.C(SESSION_COLLECTON)
	s = new(model.Session)
	err = c.Find(bson.M{"uid": user.UID}).One(&s)
	log.Println("err:%#v", s)
	if err != nil {
		// session not found. authenticate
		err = authenticateUser(db, user.UID)
		if err != nil {
			return
		}
	} else {
		// TODO: session found. update lastlogin
		log.Println("session:", s.Token)
		err = cmodel.ErrAlreadyLogin
	}
	return
}

func authenticateUser(db *mgo.Database, uid int) (err error) {
	c := db.C(AUTH_COLLECTION)
	auth := model.Auth{}
	err = c.Find(bson.M{"uid": uid}).One(&auth)
	if err != nil {
		// password not found
		fmt.Printf("password:%#v\n", auth.Password)
	} else {
		err = cmodel.ErrIncorrectPassword
		return
	}
	return
}

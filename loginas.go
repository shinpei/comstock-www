package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

const (
	AUTH_COLLECTION   string = "authinfo"
	USER_COLLECTION   string = "user"
	SESSION_COLLECTON string = "session"
)

func CheckSession(db *mgo.Database, token string) (err error) {
	c := db.C(SESSION_COLLECTON)
	q := c.Find(bson.M{"token": token})
	count, _ := q.Count()
	if count == 1 {
		// found session

	} else {
		// not found session. login requires
		log.Println("Couldn't found session")
		err = cmodel.ErrSessionNotFound
	}
	return
}

func LoginAs(db *mgo.Database, l *model.LoginRequest) (s *model.Session, err error) {
	c := db.C(USER_COLLECTION)

	item := model.User{}
	iter := c.Find(nil).Limit(100).Iter()
	for iter.Next(&item) {
		log.Printf("%#v", item)
	}

	/*
		err = c.Find(bson.M{"mail": l.Mail()}).One(&item)
		if err != nil {
			log.Println("Counln't find user")
			err = cmodel.ErrUserNotFound
			return
		}
		// found user
		log.Printf("item;%#v\n", item)
		log.Println("mail:'", item.Mail(), "'")
	*/
	return
}

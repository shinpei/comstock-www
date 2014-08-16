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
	AUTH_COLLECTION   string = "authinfo"
	USER_COLLECTION   string = "user"
	SESSION_COLLECTON string = "session"
)

func CheckSession(db *mgo.Database, authinfo string) (err error) {
	c := db.C(SESSION_COLLECTON)
	q := c.Find(bson.M{"token": authinfo})
	count, _ := q.Count()
	if count == 1 {
		// found session

	} else {
		// not found session. login requires
		fmt.Println("Couldn't found session")
		err = cmodel.ErrSessionNotFound
	}
	return
}

func LoginAs(db *mgo.Database, l *model.LoginRequest) (s *model.Session, err error) {
	c := db.C(AUTH_COLLECTION)
	iter := c.Find(nil).Iter()
	var item cmodel.AuthInfo
	count := 0
	for iter.Next(&item) {
		log.Println(item.Mail())
		count++
	}
	return
}

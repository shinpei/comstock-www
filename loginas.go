package main

import (
	"fmt"
	"github.com/shinpei/comstock-www/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

const (
	AUTH_COLLECTION   string = "authinfo"
	USER_COLLECTION   string = "user"
	SESSION_COLLECTON string = "session"
)

func CheckSession(db *mgo.Database, authinfo string) {
	c := db.C(SESSION_COLLECTON)
	q := c.Find(bson.M{"token": authinfo})
	count, _ := q.Count()
	println("count:", count)

	fmt.Printf("Query:%#v\n", q)
	return
}

func LoginAs(db *mgo.Database, l *model.LoginInfo) (res string, err error) {
	c := db.C(AUTH_COLLECTION)
	iter := c.Find(nil).Iter()
	var result model.LoginInfo
	count := 0
	for iter.Next(&result) {
		log.Println(result.Mail())
		count++
	}
	return
}

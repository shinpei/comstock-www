package main

import (
	"github.com/shinpei/comstock-www/model"
	"labix.org/v2/mgo"
	"log"
)

const (
	AUTH_COLLECTION string = "authinfo"
)

func LoginAs(db *mgo.Database, info *model.LoginInfo) {
	c := db.C(AUTH_COLLECTION)
	iter := c.Find(nil).Iter()
	var result model.LoginInfo
	for iter.Next(&result) {
		log.Println(result.Mail())
	}
}

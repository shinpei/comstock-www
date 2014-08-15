package main

import (
	"github.com/shinpei/comstock-www/model"
	"labix.org/v2/mgo"
	"log"
)

const (
	AUTH_COLLECTION string = "command"
)

func LoginAs(db *mgo.Database, info *model.LoginInfo) {
	c := db.C(AUTH_COLLECTION)
	iter := c.Find(nil).Iter()
	var result model.LoginInfo
	count := 0
	for iter.Next(&result) {
		log.Println(result.Mail())
		count++
	}
	println(count)
}
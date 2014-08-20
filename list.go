package main

import (
	"github.com/shinpei/comstock-www/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

func ListCommands(db *mgo.Database, token string) (err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}

	c := db.C(COMMAND_COLLECTION)
	cmd := model.CommandItem{}
	iter := c.Find(bson.M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()

	if err != nil {
		return
	}
	for iter.Next(&cmd) {
		log.Println("hi")
	}

	return
}

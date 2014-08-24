package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

func FetchCommandFromNumber(db *mgo.Database, token string, num int) (cmds []cmodel.Command, err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}
	// TODO: check session expiration
	c := db.C(COMMAND_COLLECTION)
	iter := c.Find(bson.M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()
	ci := model.CommandItem{}
	counter := 0
	for iter.Next(&ci) {
		counter++
		if counter == num {
			cmds = append(cmds, cmodel.Command{Cmd: ci.Data.Command, Timestamp: ci.Date})
			log.Println("Found! ", ci.Data.Command)
		}
	}
	if counter < num {
		err = cmodel.ErrCommandNotFound
		return
	}
	return
}

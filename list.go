package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

func ListCommands(db *mgo.Database, token string) (cmds []cmodel.Command, err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		if err == cmodel.ErrSessionNotFound {
			return
		} else {
			log.Fatal("Crash")
		}
	}

	c := db.C(COMMAND_COLLECTION)
	cmd := model.CommandItem{}
	iter := c.Find(bson.M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()
	for iter.Next(&cmd) {
		cmds = append(cmds, cmodel.Command{Cmd: cmd.Data.Command, Timestamp: cmd.Date})
	}
	return
}

package main

import (
	"github.com/shinpei/comstock-www/model"
	"labix.org/v2/mgo"
)

func PostCommand(db *mgo.Database, token string, cmd string) (err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}
	//TODO: check session expiration
	c := db.C(COMMAND_COLLECTION)
	command := model.CreateCommandItem(user.UID, cmd)
	err = c.Insert(command)
	return
}

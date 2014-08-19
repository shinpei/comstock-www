package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func ListCommands(db *mgo.Database, token string) (err error) {
	c := db.C(SESSION_COLLECTON)
	user := model.Session{}
	err = c.Find(bson.M{"token": token}).One(&user)
	if err != nil {
		// session not found. reject.
		err = cmodel.ErrSessionNotFound
		return
	}
	return
}

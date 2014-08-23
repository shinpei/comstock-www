package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

func RegisterUser(db *mgo.Database, mail string, password string) (err error) {
	c := db.C(USER_COLLECTION)
	user := model.User{}
	err = c.Find(bson.M{"mail": mail}).One(&user)
	if err == nil {
		// existing user.
		log.Println("Register request issued, but user ", mail, "already exist")
		err = cmodel.ErrUserAlreadyExist
		return
	}
	return
}

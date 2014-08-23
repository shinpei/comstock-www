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
	count, err := c.Find(bson.M{}).Count()
	// TODO: validate mail, password
	uid := count + 1 // TODO: also, validate uid
	newUser := model.CreateUserForNewCommer(mail, uid)
	err = c.Insert(newUser)
	if err != nil {
		err = cmodel.ErrServerSystem
		return
	}
	c = db.C(AUTH_COLLECTION)
	auth := model.CreateAuthForNewComer(uid, password)
	err = c.Insert(auth)
	if err != nil {
		err = cmodel.ErrServerSystem
		return
	}

	return
}

package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/url"
)

func RegistUserHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()

	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["mail"] == nil || m["password"] == nil {
		http.Error(w, "Invalid register request", http.StatusBadRequest)
		return
	}
	err := RegisterUser(db, m["mail"][0], m["password"][0])
	if err == cmodel.ErrUserAlreadyExist {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte("User added, thank you for registering"))

}
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

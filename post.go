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

func PostCommandHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()
	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["authinfo"] == nil || m["cmd"] == nil {
		http.Error(w, "Invalid post command requst", http.StatusBadRequest)
	}
	err := postCommand(db, m["authinfo"][0], m["cmd"][0])
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Write([]byte("Success"))
}

func postCommand(db *mgo.Database, token string, cmd string) (err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}
	//TODO: check session expiration
	c := db.C(COMMAND_COLLECTION)
	command := model.CreateCommandItem(user.UID, cmd)

	// check if there're hash commit
	ci := model.CommandItem{}
	err = c.Find(bson.M{"hash": command.Hash}).One(&ci)
	if err != nil {
		// Not Found also comes here
		log.Printf("Hash: %x\n", command.Hash)
	} else {
		log.Printf("Duplication? #%v\n", ci)
		if command.Data.Command == cmd {
			log.Printf("Duplicated: %s\n", cmd)
			return
		}
	}
	err = c.Insert(command)
	if err != nil {
		log.Println("Cannot save command", err.Error())
	}
	return
}

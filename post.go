package main

import (
	"encoding/json"
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

	// new from 0.2.0, convert recieved command json data
	ccmd := cmodel.Command{}
	err := json.Unmarshal([]byte(m["cmd"][0]), &ccmd)
	if err != nil {
		panic(err)
	}

	// actual save to the mongo
	err = postCommand(db, m["authinfo"][0], &ccmd)
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Write([]byte("Success"))
}

func postCommand(db *mgo.Database, token string, cmd *cmodel.Command) (err error) {

	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}
	//TODO: check session expiration
	c := db.C(COMMAND_COLLECTION)
	command := model.CreateCommandItem(user.UID, cmd.Cmd)

	// check if there're hash commit
	ci := model.CommandItem{}
	err = c.Find(bson.M{"hash": command.Hash}).One(&ci)
	if err != nil {
		// Not Found also comes here
		log.Printf("Hash: %x\n", command.Hash)
	} else {
		log.Printf("Duplication? #%v\n", ci)
		if command.Data.Command == cmd.Cmd {
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

/* Receive chunk data (array of commands) at the same time
* the data should be []CommnadData
 */

func PostChunkCommandsHandler(w http.ResponseWriter, req *http.Request) {

	session, db := getSessionAndDB()
	defer session.Close()
	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["authinfo"] == nil || m["cmd"] == nil {
		http.Error(w, "Invalid request for postChunkCommands", http.StatusBadRequest)
	}
	err := postChunkCommands(db, m["authinfo"][0], m["cmd"][0])
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Write([]byte("Success"))
}

func postChunkCommands(db *mgo.Database, token string, cmd string) (err error) {

	return
}

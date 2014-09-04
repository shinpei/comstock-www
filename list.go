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

func ListHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()

	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["authinfo"] == nil {
		// error
		log.Println("Error, check session requires param")
		http.Error(w, "Session check needs parameters", http.StatusBadGateway)
		return
	}
	cmds, err := ListCommands(db, m["authinfo"][0])
	if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resJson, err := json.Marshal(cmds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(resJson)
}
func ListCommands(db *mgo.Database, token string) (cmds []cmodel.Command, err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		if _, ok := err.(*cmodel.SessionNotFoundError); ok {
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

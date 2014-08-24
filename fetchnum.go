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
	"strconv"
)

func FetchHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()
	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["authinfo"] == nil || m["number"] == nil {
		http.Error(w, "Invalid fetch command request", http.StatusBadRequest)
		return
	}
	cmdNum, err := strconv.Atoi(m["number"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cmds, err := fetchCommandFromNumber(db, m["authinfo"][0], cmdNum)
	if err == cmodel.ErrSessionExpires || err == cmodel.ErrSessionNotFound {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if err == cmodel.ErrCommandNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resJson, err := json.Marshal(cmds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(string(resJson))
	w.Header().Set("Content-type", "application/json")
	w.Write(resJson)

}

func fetchCommandFromNumber(db *mgo.Database, token string, num int) (cmds []cmodel.Command, err error) {
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

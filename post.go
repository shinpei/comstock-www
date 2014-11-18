package main

import (
	"encoding/json"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"net/url"
)

func PostHistoryHandler(w http.ResponseWriter, req *http.Request) {

	session, db := getSessionAndDB()
	defer session.Close()
	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["token"] == nil || m["history"] == nil {
		http.Error(w, "Invalid post command requst", http.StatusBadRequest)
	}

	// new from 0.2.0, convert recieved command json data
	nh := cmodel.NaiveHistory{}
	err := json.Unmarshal([]byte(m["history"][0]), &nh)
	if err != nil {
		panic(err)
	}
	// actual save to the mongo
	err = postHistory(db, m["token"][0], &nh)
	//	err = postHistory(db, m["authinfo"][0],ccmd)
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Write([]byte("Success"))
}

func postHistory(db *mgo.Database, tk string, nh *cmodel.NaiveHistory) (err error) {
	usession, err := GetUserSession(db, tk)
	hist := model.TranslateNaiveHistoryToHistory(usession.UID, nh)
	err = InsertHistory(db, hist)
	if err != nil {
		log.Println("Cannot save history", err.Error())
	}
	return
}

/*
func postHistory(db *mgo.Database, tk string, cmd string, date time.Time, desc string) (err error) {
	user, err := GetUserSession(db, tk)
	hist := model.CreateHistory(user.UID, cmd, date, desc)

	err = InsertHistory(db, hist)
	if err != nil {
		log.Println("Cannot save command", err.Error())
	}
	return
}
*/

/*
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
*/

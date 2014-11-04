package main

import (
	"encoding/json"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
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
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.CommandNotFoundError); ok {
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
	return
}

func fetchHistoryFromNumber(db *mgo.Database, token string, num int) (hist *model.History, err error) {
	hist, err = findHistory(db, token, num)
	return
}

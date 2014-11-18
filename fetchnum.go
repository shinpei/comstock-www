package main

import (
	"encoding/json"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
	"strconv"
)

func FetchHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()
	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["token"] == nil || m["number"] == nil {
		http.Error(w, "Invalid fetch command request", http.StatusBadRequest)
		return
	}
	cmdNum, err := strconv.Atoi(m["number"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hist, err := fetchHistoryFromNumber(db, m["token"][0], cmdNum)
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.CommandNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if _, ok := err.(*cmodel.IllegalArgumentError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resJson, err := json.Marshal(hist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(resJson)

}
func fetchHistoryFromNumber(db *mgo.Database, tk string, num int) (nh *cmodel.NaiveHistory, err error) {

	hist, err := FindHistoryFromNum(db, tk, num)
	if err != nil {
		return
	}
	nh = model.TranslateHistoryToNaiveHistory(hist.UID, hist)
	return
}

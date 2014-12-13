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

func ListHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()

	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["token"] == nil {
		// error
		log.Println("Error, check session requires param")
		http.Error(w, "Session check needs parameters", http.StatusBadGateway)
		return
	}
	hists, err := listHistories(db, m["token"][0])
	if err != nil {
		if _, ok := err.(*cmodel.SessionNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if _, ok := err.(*cmodel.CommandNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if _, ok := err.(*cmodel.IllegalArgumentError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	resJson, err := json.Marshal(hists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(resJson)
}
func listHistories(db *mgo.Database, tk string) (nhs []*cmodel.NaiveHistory, err error) {

	//TODO: duplcation for usession
	usession, err := GetUserSession(db, tk)
	if err != nil {
		if _, ok := err.(*cmodel.SessionNotFoundError); ok {
			return
		} else {
			log.Fatal("Crash:", err.Error())
		}
	}

	hists, err := FindHistoryLastN(db, tk, 100) // defulat limit is 100
	if err != nil {
		if _, ok := err.(*cmodel.CommandNotFoundError); ok {
			return
		} else if _, ok := err.(*cmodel.IllegalArgumentError); ok {
			return
		} else {
			log.Fatal("Crash", err.Error())
		}
	}
	for _, hist := range hists {
		nhs = append(nhs, model.TranslateHistoryToNaiveHistory(usession.UID, hist))
	}
	if err != nil {
		log.Fatalln(err)
	}
	return
}

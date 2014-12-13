package main

import (
	cmodel "github.com/shinpei/comstock/model"
	"log"
	"net/http"
	"net/url"
)

func CheckSessionHandler(w http.ResponseWriter, req *http.Request) {

	session, db := getSessionAndDB()
	defer session.Close()
	// make sure param exists
	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["token"] == nil {
		// error
		log.Println("Error, check session requires token")
		http.Error(w, "session check needs parameters", http.StatusBadRequest)
		return
	}
	_, err := GetUserSession(db, m["token"][0])
	if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// if pass reaches here, session found. do nothing.
}

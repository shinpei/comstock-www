package main

import (
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
)

func TranslateHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()
	param, _ := url.ParseQuery(req.URL.RawQuery)
	if param["authinfo"] == nil || param["authinfo"] == nil {
		http.Error(w, "Invalid post command request", http.StatusBadRequest)
	}
	err := translateCommand(db, param["authinfo"][0], param["cmd"][0])
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Write([]byte("Success"))
}

func translateCommand(db *mgo.Database, token string, cmd string) (err error) {
	return
}
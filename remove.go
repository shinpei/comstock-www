package main

import (
	"labix.org/v2/mgo"
	//	"labix.org/v2/mgo/bson"
	//	"log"
	"net/http"
	"net/url"
)

func RemoveOneHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()

	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["token"] == nil || m["num"] == nil {
		http.Error(w, "Invalid removeOne request", http.StatusBadRequest)
		return
	}
	_ = db

}
func removeOne(db *mgo.Database, token string, num int) (err error) {

	return
}

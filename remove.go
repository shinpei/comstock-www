package main

import (
	//	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	//	"log"
	"net/http"
	"net/url"
	"strconv"
)

func RemoveOneHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()

	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["token"] == nil || m["index"] == nil {
		http.Error(w, "Invalid removeOne request", http.StatusBadRequest)
		return
	}
	cmdIdx, err := strconv.Atoi(m["index"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cmdIdx < 1 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = removeOne(db, m["token"][0], cmdIdx)
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if _, ok := err.(*cmodel.ServerSystemError); ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if _, ok := err.(*cmodel.CommandNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func removeOne(db *mgo.Database, tk string, idx int) (err error) {
	user, err := GetUserSession(db, tk)
	if err != nil {
		return
	}
	err = RemoveHistoryNth(db, user.UID, idx)
	return
}

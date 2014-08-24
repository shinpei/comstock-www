package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"net/http"
	"net/url"
)

func PostCommandHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()
	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["authinfo"] == nil || m["cmd"] == nil {
		http.Error(w, "Invalid post command requst", http.StatusBadRequest)
	}
	err := postCommand(db, m["authinfo"][0], m["cmd"][0])
	if err == cmodel.ErrSessionExpires || err == cmodel.ErrSessionNotFound {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Write([]byte("Success"))
}

func postCommand(db *mgo.Database, token string, cmd string) (err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}
	//TODO: check session expiration
	c := db.C(COMMAND_COLLECTION)
	command := model.CreateCommandItem(user.UID, cmd)
	err = c.Insert(command)
	return
}

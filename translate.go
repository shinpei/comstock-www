package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
)

func RanslateHandler(w http.ResponseWriter, req *http.Request) {
	session, db := getSessionAndDB()
	defer session.Close()
	param, _ := url.ParseQuery(req.URL.RawQuery)
	if param["authinfo"] == nil || param["authinfo"] == nil {
		http.Error(w, "Invalid post command request", http.StatusBadRequest)
		return
	}
	err := translateCommand(db, param["authinfo"][0])
	if _, ok := err.(*cmodel.SessionExpiresError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if _, ok := err.(*cmodel.SessionNotFoundError); ok {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Write([]byte("Success"))
}

func translateCommand(db *mgo.Database, token string) (err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}
	c := db.C(COMMAND_COLLECTION)
	iter := c.Find(bson.M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()
	ci := model.OldCommandItem{}
	for iter.Next(&ci) {
		hist := model.TranslateCommand1to2(&ci)
		println(hist.Command())

	}

	return
}

// FOR DEBUG
package main

import (
	"github.com/shinpei/comstock-www/model"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/url"
)

func TransHandler(w http.ResponseWriter, req *http.Request) {

	session, db := getSessionAndDB()
	defer session.Close()
	m, _ := url.ParseQuery(req.URL.RawQuery)
	if m["authinfo"] == nil {
		http.Error(w, "Invalid trans command request", http.StatusBadRequest)
	}
	token := m["authinfo"][0]

	user, err := GetUserSession(db, token)
	if err != nil {
		panic(err)
	}
	c := db.C(COMMAND_COLLECTION)
	iter := c.Find(bson.M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()
	ci := model.CommandItem{}
	for iter.Next(&ci) {
		log.Println("data==>", ci.Data.Command)
		hist := model.TranslateCommand1to2(&ci)
		log.Println("hist==>", hist.FlowPtr.ItemsPtr[0].Command)
	}
	return
}

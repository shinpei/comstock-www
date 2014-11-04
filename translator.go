// FOR DEBUG
package main

import (
	"github.com/shinpei/comstock-www/model"
	"net/http"
	"net/url"
)

type M map[string]interface{}

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
	c := NewCommandStorage(db)
	iter := c.FindFor(M{"uid": user.UID}, 100)
	defer iter.Close()
	ci := model.OldCommandItem{}
	for iter.Next(&ci) {
		hist := model.TranslateCommand1to2(&ci)
		InsertHistory(db, hist)
	}
	return
}

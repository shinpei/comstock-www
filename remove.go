package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
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
	err = removeOne(db, m["token"][0], cmdIdx)
	if err == cmodel.ErrSessionExpires || err == cmodel.ErrSessionNotFound {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	} else if err == cmodel.ErrServerSystem {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == cmodel.ErrCommandNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func removeOne(db *mgo.Database, token string, num int) (err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}
	c := db.C(COMMAND_COLLECTION)
	iter := c.Find(bson.M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()
	counter := 0
	ci := model.CommandItem{}
	for iter.Next(&ci) {
		counter++
		if counter == num {
			id := ci.ID

			err = c.RemoveId(id)
			if err != nil {
				log.Println("Cannot delete")
				err = cmodel.ErrServerSystem
				return
			}
		}
	}
	if counter < num {
		err = cmodel.ErrCommandNotFound
		return
	}

	return
}

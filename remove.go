package main

import (
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
	"labix.org/v2/mgo"
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

func removeOne(db *mgo.Database, token string, num int) (err error) {
	user, err := GetUserSession(db, token)
	if err != nil {
		return
	}
	c := db.C(COMMAND_COLLECTION)
	iter := c.Find(M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()
	counter := 0
	ci := model.OldCommandItem{}
	for iter.Next(&ci) {
		counter++
		if counter == num {
			id := ci.ID

			err = c.RemoveId(id)
			if err != nil {
				log.Println("Cannot delete")
				err = &cmodel.ServerSystemError{}
				return
			}
		}
	}
	if counter < num {
		err = &cmodel.CommandNotFoundError{}
		return
	}

	return
}

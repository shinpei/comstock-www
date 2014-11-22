// FOR DEBUG
package main

import (
	"fmt"
	"github.com/shinpei/comstock-www/model"
	"net/http"
	//	"net/url"
)

type M map[string]interface{}

func TransHandler(w http.ResponseWriter, req *http.Request) {

	session, db := getSessionAndDB()
	defer session.Close()
	// DANGER:
	c := db.C(USER_COLLECTION)
	iter := c.Find(M{}).Iter()
	defer iter.Close()
	u := model.User{}
	for iter.Next(&u) {
		fmt.Printf("user:%#v\n", u)
		uid := u.UID
		c2 := db.C(COMMAND_COLLECTION)
		iter := c2.Find(M{"uid": uid}).Iter()
		cmd := model.OldCommandItem{}
		for iter.Next(&cmd) {
			fmt.Printf("cmmand:%#v\n", cmd)
			hist := model.TranslateCommand1to2(&cmd)
			InsertHistory(db, hist)
		}
	}
	/*	user, err := GetUserSession(db, token)
		if err != nil {
			panic(err)
		}
		c := db.C(COMMAND_COLLECTION)
		iter := c.Find(M{"uid": user.UID}).Limit(100).Iter()
		defer iter.Close()
		ci := model.OldCommandItem{}
		for iter.Next(&ci) {
			hist := model.TranslateCommand1to2(&ci)
			InsertHistory(db, hist)
		}
	*/
	return
}

package main

import (
	"github.com/shinpei/comstock-www/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	MONGO_DATABSE_NAME string = "comstock-www"
	AUTH_COLLECTION    string = "authinfo"
	USER_COLLECTION    string = "user"
	SESSION_COLLECTON  string = "session"
	COMMAND_COLLECTION string = "commands"
	FLOW_COLLECTION    string = "flow"
	HISTORY_COLLECTION string = "history"
)

func getSessionAndDB() (*mgo.Session, *mgo.Database) {
	mongoURI := os.Getenv("MONGOHQ_URL")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost/" + MONGO_DATABSE_NAME
	}
	u, err := url.Parse(mongoURI)
	if err != nil {
		panic("couldn't parse mongouri")
	}
	dbname := u.Path
	dbname = dbname[1:] // remove slash
	session, err := mgo.DialWithTimeout(mongoURI, time.Duration(3)*time.Second)
	if err != nil {
		panic("Coulnd't dial")
	}
	session.SetSafe(&mgo.Safe{})
	return session, session.DB(dbname)
}

func InsertCommandItem(db *mgo.Database, cmd *model.NewCommandItem) {
	c := db.C(COMMAND_COLLECTION)
	ci := model.NewCommandItem{}
	err := c.Find(bson.M{"hash": cmd.Hash}).One(&ci)
	if err == nil {
		log.Printf("Duplicated? %s\n", cmd.Command)
		if cmd.Command == ci.Command {
			log.Println("Duplicated!!", cmd.Command)
			// TODO: need to count up?
			return
		}
	}
	err = c.Insert(cmd)
	if err != nil {
		log.Printf("Cannot save command, %#v, %#v\n", cmd, err)
	}
}

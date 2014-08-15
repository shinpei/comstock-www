package main

import (
	"labix.org/v2/mgo"
	"net/url"
	"os"
	"time"
)

const (
	MONGO_DATABSE_NAME string = "comstock-www"
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
	println("dbname=>", dbname)
	session, err := mgo.DialWithTimeout(mongoURI, time.Duration(3)*time.Second)
	if err != nil {
		panic("Coulnd't dial")
	}
	session.SetSafe(&mgo.Safe{})
	return session, session.DB(dbname)
}

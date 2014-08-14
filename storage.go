package main

import (
	"labix.org/v2/mgo"
	"os"
	"time"
)

const (
	MONGO_DATABSE_NAME    string = "comstock"
	MONGO_COLLECTION_USER        = "user"
)

func getSessionAndDB() (*mgo.Session, *mgo.Database) {
	mongoURI := os.Getenv("MONGOHQ_URL")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost/" + MONGO_DATABSE_NAME
	}
	session, err := mgo.DialWithTimeout(mongoURI, time.Duration(3)*time.Second)
	if err != nil {
		panic("Coulnd't dial")
	}
	session.SetSafe(&mgo.Safe{})
	return session, session.DB(MONGO_DATABSE_NAME)
}

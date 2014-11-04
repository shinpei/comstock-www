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

func InsertCommandItem(db *mgo.Database, cmd *model.NewCommandItem) (err error) {
	c := db.C(COMMAND_COLLECTION)
	ci := model.NewCommandItem{}
	err = c.Find(M{"hash": cmd.Hash}).One(&ci)
	if err == nil {
		if cmd.Command == ci.Command {
			log.Printf("Duplicated! '%s'\n", cmd.Command)
			// TODO: need to count up?

			return
		}
	}
	err = c.Insert(cmd)
	if err != nil {
		log.Printf("Cannot save command, %#v, %#v\n", cmd, err)
	}
	return
}

type history struct {
	UID         int
	Date        time.Time
	Description string
	Flow        bson.ObjectId
}

func encodeHistory(h *model.History) *history {
	return &history{
		UID:         h.UID,
		Date:        h.Date,
		Description: h.Description,
		Flow:        h.Flow,
	}
}

func decodeHistory(h *history) *model.History {
	// decoding require additional read from database, for
	// reconstructing flow, command items data structure
	return nil
}

func InsertHistory(db *mgo.Database, hist *model.History) (err error) {
	c := db.C(HISTORY_COLLECTION)

	histBuf := &history{}
	err = c.Find(M{"date": hist.Date}).One(&histBuf)
	if err == nil {
		log.Printf("Duplicated timestamp, %#x\n", hist.Date)
		return
	}

	err = insertFlow(db, hist.FlowPtr)
	if err != nil {
		return
	}
	// remove FlowPtr
	h := encodeHistory(hist)
	err = c.Insert(h)
	if err != nil {
		log.Println("Cannot save history", err.Error())
	}

	return
}

type flow struct {
	ID    bson.ObjectId
	Items []bson.ObjectId
}

func encodeFlow(mf *model.Flow) *flow {
	return &flow{
		ID:    mf.ID,
		Items: mf.Items,
	}
}
func decodeFlow(f *flow) *model.Flow {

	return nil
}

// This is private, we cannot use this from external
func insertFlow(db *mgo.Database, mf *model.Flow) (err error) {
	c := db.C(FLOW_COLLECTION)

	for idx, ci := range mf.ItemsPtr {
		_ = idx
		err = InsertCommandItem(db, ci)
		if err != nil {
			return
		}
	}
	f := encodeFlow(mf)
	err = c.Insert(f)
	if err != nil {
		log.Println("Cannot save flow", err.Error())
	}
	return
}

func findHistory(db *mgo.Database, token string, num int) (hist *model.History, err error) {
	user, err := GetUserSession(db, token)
	c := db.C(HISTORY_COLLECTION)
	iter := c.Find(M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()
	h := history{}
	counter := 0
	for iter.Next(&h) {
		counter++
		if counter == num {
			hist = decodeHistory(&h)
			break
		}
	}
	return
}

package main

import (
	"fmt"
	"github.com/shinpei/comstock-www/model"
	cmodel "github.com/shinpei/comstock/model"
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

//================= structs =========================
type history struct {
	UID         int
	Date        time.Time
	Description string
	Flow        bson.ObjectId
}

type flow struct {
	ID    bson.ObjectId
	Items []bson.ObjectId
}

//================= converters =========================
func encodeHistory(h *model.History) *history {
	return &history{
		UID:         h.UID,
		Date:        h.Date,
		Description: h.Description,
		Flow:        h.Flow,
	}
}

// decoding require additional read from database, for
// reconstructing flow, command items data structure
func decodeHistory(h *history, mf *model.Flow) *model.History {

	return &model.History{
		UID:         h.UID,
		Date:        h.Date,
		Description: h.Description,
		Flow:        h.Flow,
		FlowPtr:     mf,
	}
}

func encodeFlow(mf *model.Flow) *flow {
	return &flow{
		ID:    mf.ID,
		Items: mf.Items,
	}
}

func decodeFlow(f *flow, mcis []*model.NewCommandItem) *model.Flow {

	return &model.Flow{
		ID:       f.ID,
		Items:    f.Items,
		ItemsPtr: mcis,
	}
	return nil
}

//================= functions =========================
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

func FindHistoryFromNum(db *mgo.Database, tk string, num int) (hist *model.History, err error) {
	if num < 1 {
		// TODO: IllegalArgumentError()?
		err = &cmodel.CommandNotFoundError{}
		return
	}
	user, err := GetUserSession(db, tk)
	c := db.C(HISTORY_COLLECTION)
	iter := c.Find(M{"uid": user.UID}).Limit(100).Iter()
	defer iter.Close()
	h := history{}
	counter := 0
	for iter.Next(&h) {
		counter++
		if counter == num {
			// findFlow
			D("HI")
			mf, err := findFlow(db, h.Flow)
			if err != nil {
				println("Flow is nil")
				break
			}
			hist = decodeHistory(&h, mf)

			break
		}
	}
	println("Num : ", num, ", counter:", counter)
	if counter != num || hist == nil {
		err = &cmodel.CommandNotFoundError{}
	}

	return
}

// Shouldn't call from external package
func findFlow(db *mgo.Database, fID bson.ObjectId) (mf *model.Flow, err error) {

	c := db.C(FLOW_COLLECTION)
	f := flow{}
	D("Finding flow: %v\n", fID)
	err = c.Find(M{"id": fID}).One(&f)
	if err != nil {
		fmt.Println("Not found:", fID)
		return
	}
	// find CommandItems
	cis, err := findCommandItems(db, f.Items)
	mf = decodeFlow(&f, cis)
	return
}

func findCommandItems(db *mgo.Database, cIDs []bson.ObjectId) (mcis []*model.NewCommandItem, err error) {
	c := db.C(COMMAND_COLLECTION)
	for _, cid := range cIDs {
		cmd := model.NewCommandItem{}
		err = c.Find(M{"id": cid}).One(&cmd)
		if err != nil {
			// cannot find, database corrupy?
			break
		}
		mcis = append(mcis, &cmd)
	}
	return
}

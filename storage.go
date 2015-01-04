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
	AUTH_COLLECTION    string = "authinfo" // have to rename...
	USER_COLLECTION    string = "user"
	SESSION_COLLECTION string = "session"
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
type historyWithId struct {
	ID   bson.ObjectId `bson:"_id"`
	UID  int
	Flow bson.ObjectId
}

type flow struct {
	ID    bson.ObjectId
	Items []model.CommandId
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

func decodeFlow(f *flow, mcis []*model.CommandItem) *model.Flow {

	return &model.Flow{
		ID:       f.ID,
		Items:    f.Items,
		ItemsPtr: mcis,
	}
}

//================= functions =========================
func getSessionAndDB() (*mgo.Session, *mgo.Database) {
	mongoURI := os.Getenv("MONGOHQ_URL")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost/" + MONGO_DATABSE_NAME
	}
	u, err := url.Parse(mongoURI)
	if err != nil {
		log.Fatalln("Couldn't parse mongouri", mongoURI)
	}
	dbname := u.Path
	dbname = dbname[1:] // remove slash
	session, err := mgo.DialWithTimeout(mongoURI, time.Duration(3)*time.Second)
	if err != nil {
		log.Fatalln("Coulnd't dial to mongo server")
	}
	session.SetSafe(&mgo.Safe{})
	return session, session.DB(dbname)
}

func InsertCommandItem(db *mgo.Database, cmd *model.CommandItem) (err error) {

	c := db.C(COMMAND_COLLECTION)
	ci := model.CommandItem{}
	err = c.Find(M{"hash": cmd.Hash, "hitcount": 1}).One(&ci)
	if err == nil {
		// err == nil means, we found
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
	histBuf := histPool.Get().(*history)
	err = c.Find(M{"date": hist.Date}).One(histBuf)
	if err == nil {
		//TODO: see UID too
		log.Printf("Duplicated timestamp, %#x\n", hist.Date)
		return
	}
	histPool.Put(histBuf)

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
	for _, ci := range mf.ItemsPtr {
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

func FindHistoryLastN(db *mgo.Database, tk string, limit int) (hists []*model.History, err error) {
	if limit < 1 {
		err = &cmodel.IllegalArgumentError{}
		return
	}
	user, err := GetUserSession(db, tk)
	c := db.C(HISTORY_COLLECTION)
	q := c.Find(M{"uid": user.UID})
	count, _ := q.Count()
	if limit < count {
		count = limit
	}

	iter := q.Limit(count).Iter()
	defer iter.Close()
	hists = make([]*model.History, 0, count)
	h := histPool.Get().(*history)
	for iter.Next(h) {
		mf, err := findFlow(db, h.Flow)
		if err != nil {
			log.Println("Flow is nil")
			break
		}
		hist := decodeHistory(h, mf)
		hists = append(hists, hist)
	}
	return
}

func FindHistoryFromNum(db *mgo.Database, tk string, num int) (hist *model.History, err error) {
	if num < 1 {
		err = &cmodel.IllegalArgumentError{}
		return
	}
	user, err := GetUserSession(db, tk)
	c := db.C(HISTORY_COLLECTION)
	q := c.Find(M{"uid": user.UID})
	count, _ := q.Count()
	if count < num {
		err = &cmodel.CommandNotFoundError{}
		return
	}
	// default check
	if 100 < count {
		count = 100
	}
	iter := q.Limit(count).Iter()

	defer iter.Close()
	h := histPool.Get().(*history)
	counter := 0
	for iter.Next(h) {
		counter++
		if counter == num {
			// findFlow
			mf, err := findFlow(db, h.Flow)
			if err != nil {
				log.Println("Flow is nil")
				break
			}
			hist = decodeHistory(h, mf)
			break
		}
	}
	histPool.Put(h)
	if counter != num || hist == nil {
		err = &cmodel.IllegalArgumentError{}
	}

	return
}

// Shouldn't call from external package
func findFlow(db *mgo.Database, fID bson.ObjectId) (mf *model.Flow, err error) {

	c := db.C(FLOW_COLLECTION)
	f := flowPool.Get().(*flow)
	err = c.Find(M{"id": fID}).One(f)
	if err != nil {
		fmt.Println("Not found:", fID)
		return
	}
	// find CommandItems
	cis, err := findCommandItems(db, f.Items)
	if err != nil {
		fmt.Println("Command not found: ", f.Items)
	}
	mf = decodeFlow(f, cis)
	flowPool.Put(f)
	return
}

func findCommandItems(db *mgo.Database, cIDs []model.CommandId) (mcis []*model.CommandItem, err error) {
	c := db.C(COMMAND_COLLECTION)
	for _, cid := range cIDs {
		//		cmd := model.CommandItem{}
		cmd := &model.CommandItem{}
		err = c.Find(M{"hash": cid, "hitcount": 1}).One(cmd) // FIXME: for olddata
		if err != nil {
			// cannot find, database corrupt?
			break
		}
		mcis = append(mcis, cmd)
	}
	return
}

func removeFlow(db *mgo.Database, fID bson.ObjectId) (err error) {

	c := db.C(FLOW_COLLECTION)
	f := flowPool.Get().(*flow)
	err = c.Find(M{"id": fID}).One(f)
	if err != nil {
		fmt.Println("not fount for flow id=", fID)
		return
	}
	flowPool.Put(f)
	// TODO : delete command item??
	err = c.Remove(M{"id": fID})
	if err != nil {
		if err != mgo.ErrNotFound {
			fmt.Println("Cannot delete")
			err = &cmodel.ServerSystemError{}
		}
	}
	return
}

func RemoveHistoryNth(db *mgo.Database, uid int, idx int) (err error) {

	c := db.C(HISTORY_COLLECTION)
	q := c.Find(M{"uid": uid})
	count, err := q.Count()
	if count < idx {
		err = &cmodel.CommandNotFoundError{}
	}
	if 100 < count {
		count = 100
	}
	iter := q.Limit(count).Iter()
	defer iter.Close()
	counter := 0
	hwi := &historyWithId{}
	for iter.Next(hwi) {
		counter++
		if counter == idx {
			err = removeFlow(db, hwi.Flow)
			if err != nil {
				return
			}
			id := hwi.ID
			err = c.RemoveId(id)
			if err != nil {
				log.Println("Cannot delete history")
				err = &cmodel.ServerSystemError{}
				return
			}
			break
		}
	}
	if counter < idx {
		err = &cmodel.CommandNotFoundError{}
		return
	}
	return

}

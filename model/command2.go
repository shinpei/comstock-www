package model

import (
	"crypto/sha1"
	"io"
	"labix.org/v2/mgo/bson"
	"time"
)

type History struct {
	UID         int
	Date        time.Time
	Description string
	Flow        bson.ObjectId
}

type Flow struct {
	ID    bson.ObjectId
	Items []bson.ObjectId
}

type NewCommandItem struct {
	ID      bson.ObjectId
	Hash    []byte
	Command string
	Count   int
}

func CreateNewCommandData(uid int, cmd string, desc string) *History {
	h := sha1.New()
	io.WriteString(h, cmd)
	ciID := bson.NewObjectId()
	ci := &NewCommandItem{
		ID:      ciID,
		Hash:    h.Sum(nil),
		Command: cmd,
		Count:   1,
	}
	_ = ci
	fID := bson.NewObjectId()
	f := &Flow{
		ID: fID,
		//		Items: []bson.ObjectId{ciID},
	}
	_ = f
	return &History{
		UID:         uid,
		Date:        time.Now(),
		Description: desc,
		Flow:        fID,
	}

}

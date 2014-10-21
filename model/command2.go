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
	FlowPtr     *Flow
}

type Flow struct {
	ID       bson.ObjectId
	Items    []bson.ObjectId
	ItemsPtr []*NewCommandItem
}

type NewCommandItem struct {
	ID       bson.ObjectId
	Hash     []byte
	Command  string
	HitCount int
}

func CreateNewCommandData(uid int, cmd string, desc string) *History {
	h := sha1.New()
	io.WriteString(h, cmd)
	ciID := bson.NewObjectId()
	ci := &NewCommandItem{
		ID:       ciID,
		Hash:     h.Sum(nil),
		Command:  cmd,
		HitCount: 1,
	}
	fID := bson.NewObjectId()
	f := &Flow{
		ID:       fID,
		Items:    []bson.ObjectId{ciID},
		ItemsPtr: []*NewCommandItem{ci},
	}
	return &History{
		UID:         uid,
		Date:        time.Now(),
		Description: desc,
		Flow:        fID,
		FlowPtr:     f,
	}

}

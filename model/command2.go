package model

import (
	"crypto/sha1"
	"io"
	"labix.org/v2/mgo/bson"
	"time"
)

type NewCommandItem struct {
	ID          bson.ObjectId
	UID         int
	Date        time.Time
	Hash        []byte
	Command     string
	Description string
}

func CreateNewCommandData(uid int, cmd string, desc string) *NewCommandItem {
	h := sha1.New()
	io.WriteString(h, cmd)
	return &NewCommandItem{
		ID:          bson.NewObjectId(),
		UID:         uid,
		Hash:        h.Sum(nil),
		Date:        time.Now(),
		Command:     cmd,
		Description: desc,
	}
}

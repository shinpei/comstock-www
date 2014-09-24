package model

import (
	"crypto/sha1"
	"io"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type CommandItem struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	UID  int
	Date string `json:"date" bson:"date"` // TODO: fix it with time.Time
	Hash []byte
	Data CommandData // TODO: fix name
}

type CommandData struct {
	Command string // need to
	Desc    string
}

func CreateCommandItem(uid int, cmd string) *CommandItem {
	h := sha1.New()
	io.WriteString(h, cmd)
	return &CommandItem{ID: bson.NewObjectId(), UID: uid, Hash: h.Sum(nil), Date: strconv.FormatInt(time.Now().Unix()*1000, 10), Data: CommandData{Command: cmd, Desc: ""}}
}

package model

import (
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type CommandItem struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	UID  int
	Date string      `json:"date" bson:"date"` // TODO: fix it with time.Time
	Data CommandData // TODO: fix name
}

type CommandData struct {
	Command string // need to
	Desc    string
}

func CreateCommandItem(uid int, cmd string) *CommandItem {
	return &CommandItem{ID: bson.NewObjectId(), UID: uid, Date: strconv.FormatInt(time.Now().Unix()*1000, 10), Data: CommandData{Command: cmd, Desc: ""}}
}

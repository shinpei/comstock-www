package model

import (
	"crypto/sha1"
	cmodel "github.com/shinpei/comstock/model"
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

func CreateCommandItemFromCommand(uid int, cmd *cmodel.Command) *CommandItem {
	h := sha1.New()
	io.WriteString(h, cmd.Cmd)
	return &CommandItem{ID: bson.NewObjectId(), UID: uid, Hash: h.Sum(nil), Date: strconv.FormatInt(time.Now().Unix()*1000, 10), Data: CommandData{Command: cmd.Cmd, Desc: ""}}
}

func CreateCommandItem(uid int, cmd string) *CommandItem {
	h := sha1.New()
	io.WriteString(h, cmd)
	return &CommandItem{ID: bson.NewObjectId(), UID: uid, Hash: h.Sum(nil), Date: strconv.FormatInt(time.Now().Unix()*1000, 10), Data: CommandData{Command: cmd, Desc: ""}}
}

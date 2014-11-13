package model

import (
	"crypto/sha1"
	cmodel "github.com/shinpei/comstock/model"
	"io"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type OldCommandItem struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	UID  int
	Date string `json:"date" bson:"date"` // TODO: fix it with time.Time
	Hash []byte
	Data OldCommandData // TODO: fix name
}

type OldCommandData struct {
	Command string // need to
	Desc    string
}

func CreateOleCommandItemFromCommand(uid int, cmd *cmodel.Command) *OldCommandItem {
	h := sha1.New()
	io.WriteString(h, cmd.Cmd)
	return &OldCommandItem{ID: bson.NewObjectId(), UID: uid, Hash: h.Sum(nil), Date: strconv.FormatInt(time.Now().Unix()*1000, 10), Data: OldCommandData{Command: cmd.Cmd, Desc: ""}}
}

func CreateOldCommandItem(uid int, cmd string) *OldCommandItem {
	h := sha1.New()
	io.WriteString(h, cmd)
	return &OldCommandItem{ID: bson.NewObjectId(), UID: uid, Hash: h.Sum(nil), Date: strconv.FormatInt(time.Now().Unix()*1000, 10), Data: OldCommandData{Command: cmd, Desc: ""}}
}

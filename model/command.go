package model

import (
	"crypto/sha1"
	"io"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type CommandId []byte

type History struct {
	UID         int
	Date        time.Time
	Description string
	Flow        bson.ObjectId
	FlowPtr     *Flow
}

func (h *History) Command() string {
	l := len(h.FlowPtr.ItemsPtr)
	if l == 1 {
		return h.FlowPtr.ItemsPtr[0].Command
	} else if l > 1 {
		buf := ""
		for _, item := range h.FlowPtr.ItemsPtr {
			buf += item.Command + " => "
		}
		return buf[:len(buf)-4]
	} else {
		panic("Received size=0 history as receiver")
	}
}

type Flow struct {
	ID       bson.ObjectId
	Items    []CommandId
	ItemsPtr []*CommandItem
}

func (f *Flow) Command() string {
	dummyH := &History{
		UID:         0,
		Date:        time.Now(),
		Description: "",
		Flow:        f.ID,
		FlowPtr:     f,
	}
	return dummyH.Command()
}

// use command hash as ID
// WARN: optimistic assumption -> command hash won't collide
type CommandItem struct {
	//	ID       bson.ObjectId
	Hash     CommandId
	Command  string
	HitCount int
}

// WARN: Command cannot be created freely.
func CreateHistoryFromFlow(uid int, date time.Time, desc string, f *Flow) *History {
	return &History{
		UID:         uid,
		Date:        date,
		Description: desc,
		Flow:        f.ID,
		FlowPtr:     f,
	}
}

// WARN: Command cannot be created freely.
func CreateCommandItem(cmd string) (cid CommandId, item *CommandItem) {
	// make sure uid, cmd is not nil
	h := sha1.New()
	if cmd == "" {
		return
	}
	io.WriteString(h, cmd)
	//	hash = h.Sum(nil)
	cid = h.Sum(nil) // we use hash as cid
	item = &CommandItem{
		Hash:     cid,
		Command:  cmd,
		HitCount: 1,
	}
	return
}

func CreateFlow(cis []*CommandItem) (fID bson.ObjectId, f *Flow) {
	cilen := len(cis)
	if cilen == 0 {
		// return with nil
		return
	}
	items := []CommandId{}
	for _, ci := range cis {
		items = append(items, ci.Hash)
	}
	fID = bson.NewObjectId()
	f = &Flow{
		ID:       fID,
		Items:    items,
		ItemsPtr: cis,
	}
	return
}

func CreateHistory(uid int, cmds []string, date time.Time, desc string) (his *History) {
	if len(cmds) == 1 {
		_, ci := CreateCommandItem(cmds[0])
		_, f := CreateFlow([]*CommandItem{ci})
		his = CreateHistoryFromFlow(uid, date, desc, f)
	} else if len(cmds) > 1 {
		cis := []*CommandItem{}
		for _, cmd := range cmds {
			_, ci := CreateCommandItem(cmd)
			cis = append(cis, ci)
		}
		_, f := CreateFlow(cis)
		his = CreateHistoryFromFlow(uid, date, desc, f)
	} else {
		log.Fatalln("len(cmds) is 0, no history made")
	}
	return
}

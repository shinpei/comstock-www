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

func (h *History) Command() string {
	length := len(h.FlowPtr.ItemsPtr)
	if length == 1 {
		return h.FlowPtr.ItemsPtr[0].Command
	} else if length > 1 {
		buf := ""
		for _, item := range h.FlowPtr.ItemsPtr {
			buf += item.Command + " => "
		}
		return buf[:len(buf)-4]
	} else {
		panic("Null")
	}
}

type Flow struct {
	ID       bson.ObjectId
	Items    []bson.ObjectId
	ItemsPtr []*NewCommandItem
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

type NewCommandItem struct {
	ID       bson.ObjectId
	Hash     []byte
	Command  string
	HitCount int
}

func CreateHistory(uid int, cmds []string) *History {
	histLen := len(cmds)
	if histLen == 0 {
		return nil
	} else {
		for _, cmd := range cmds {
			println(cmd)
		}
	}
	return nil
}

func CreateNewCommandItem(uid int, cmd string) (id bson.ObjectId, item *NewCommandItem) {
	// make sure uid, cmd is not nil
	h := sha1.New()
	if cmd == "" {
		return
	}
	id = bson.NewObjectId()
	io.WriteString(h, cmd)
	item = &NewCommandItem{
		ID:       id,
		Hash:     h.Sum(nil),
		Command:  cmd,
		HitCount: 1,
	}
	return
}

func CreateFlow(cis []*NewCommandItem) (fID bson.ObjectId, f *Flow) {
	cilen := len(cis)
	if cilen == 0 {
		// return with nil
		return
	}
	items := []bson.ObjectId{}
	for _, ci := range cis {
		items = append(items, ci.ID)
	}
	fID = bson.NewObjectId()
	f = &Flow{
		ID:       fID,
		Items:    items,
		ItemsPtr: cis,
	}
	return
}

func CreateNewHistory(uid int, cmd string, date time.Time, desc string) *History {
	ciID, ci := CreateNewCommandItem(uid, cmd)
	fID := bson.NewObjectId()
	f := &Flow{
		ID:       fID,
		Items:    []bson.ObjectId{ciID},
		ItemsPtr: []*NewCommandItem{ci},
	}
	return &History{
		UID:         uid,
		Date:        date,
		Description: desc,
		Flow:        fID,
		FlowPtr:     f,
	}

}

package model

import (
	"crypto/sha1"
	"io"
	"strconv"
	"time"
)

func TranslateCommand1to2(item *CommandItem) *NewCommandItem {

	h := sha1.New()
	var hash []byte
	if item.Hash == nil {
		io.WriteString(h, item.Data.Command)
		hash = h.Sum(nil)
	}
	timeVal, _ := strconv.ParseInt(item.Date, 10, 64)
	return &NewCommandItem{
		ID:          item.ID,
		UID:         item.UID,
		Hash:        hash,
		Date:        time.Unix(timeVal, 0),
		Command:     item.Data.Command,
		Description: item.Data.Desc,
	}
}

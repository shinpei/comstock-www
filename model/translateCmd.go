package model

import (
	"crypto/sha1"
	"io"
	//	"strconv"
)

func TranslateCommand1to2(item *CommandItem) *History {
	h := sha1.New()
	var hash []byte
	if item.Hash == nil {
		io.WriteString(h, item.Data.Command)
		hash = h.Sum(nil)
	}
	_ = hash
	/*
		timeVal, err := strconv.ParseInt(item.Date, 10, 64)
			if err != nil {
				panic(err)
			}
			// check if command are
			//
	*/
	return &History{}

}

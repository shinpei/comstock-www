package model

import (
	"strconv"
	"time"
)

// suppose Single item will come
func TranslateCommand1to2(item *OldCommandItem) *History {
	// parse to int64
	timeVal, err := strconv.ParseInt(item.Date, 10, 64)
	date := time.Unix(0, timeVal*1000000)
	println(date.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
	// create CommandDataStructure
	return CreateHistory(item.UID, item.Data.Command, date, item.Data.Desc)
}

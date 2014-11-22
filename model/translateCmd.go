package model

import (
	cmodel "github.com/shinpei/comstock/model"
	"log"
	"strconv"
	"time"
)

// suppose Single item will come
func TranslateCommand1to2(item *OldCommandItem) *History {
	// parse to int64
	var date time.Time
	if item.Date == "" {
		date = time.Now()
	} else {
		timeVal, err := strconv.ParseInt(item.Date, 10, 64)
		if err != nil {
			log.Println("Couldn't parse time=", date, err)
			date = time.Now()
		}
		date = time.Unix(0, timeVal*1000000)
	}
	// create CommandDataStructure
	return CreateHistory(item.UID, []string{item.Data.Command}, date, item.Data.Desc)
}

func TranslateNaiveHistoryToHistory(uid int, item *cmodel.NaiveHistory) *History {
	return CreateHistory(uid, item.Cmds, item.Date, item.Description)
}

func TranslateHistoryToNaiveHistory(uid int, hist *History) (nh *cmodel.NaiveHistory) {
	l := len(hist.FlowPtr.Items)
	if l == 1 {
		nh = &cmodel.NaiveHistory{
			Date:        hist.Date,
			Description: hist.Description,
			Cmds:        []string{hist.Command()},
			Shell:       "",
		}

	} else if l > 1 {
		var cmds []string
		for _, cmdPtr := range hist.FlowPtr.ItemsPtr {
			cmds = append(cmds, cmdPtr.Command)
		}
		nh = &cmodel.NaiveHistory{
			Date:        hist.Date,
			Description: hist.Description,
			Cmds:        cmds,
			Shell:       "",
		}

	} else {
		panic("Cannot translate")
		nh = nil
	}
	return
}

package model

import (
	"time"
)

type CommandItem struct {
	UID  int
	Date string      // TODO: fix it with time.Time
	Data CommandData // TODO: fix name
}

type CommandData struct {
	Command string // need to
	Desc    string
}

func CreateCommandItem(uid int, cmd string) *CommandItem {
	return &CommandItem{UID: uid, Date: string(time.Now().Unix()), Data: CommandData{Command: cmd, Desc: ""}}
}

package model

import (
	"time"
)

type Auth struct {
	UID      int
	Date     int64 // TODO: make it with time pkg
	Password string
}

func CreateAuthForNewComer(uid int, pass string) *Auth {
	return &Auth{UID: uid, Date: time.Now().Unix(), Password: pass}
}

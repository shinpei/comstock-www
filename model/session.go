package model

import (
	"time"
)

type Session struct {
	Token   string
	UID     int
	Expires int64 // TODO: replace it with time
}

func CreateSession(token string, uid int) *Session {
	return &Session{Token: token, UID: uid, Expires: time.Now().Unix()}
}

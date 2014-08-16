package model

import (
	"time"
)

type Session struct {
	token   string
	uid     int
	expires time.Time
}

func CreateSession(token string, uid int) *Session {
	return &Session{token: token, uid: uid}
}

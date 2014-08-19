package model

import (
//	"time"
)

type User struct {
	Mail      string
	Username  string
	UID       int
	Created   int // TODO: make it with time package
	LastLogin int // TODO: make it with time package
}

func CreateUser(mail string, username string, uid int) *User {
	return &User{Mail: mail, Username: username, UID: uid}
}

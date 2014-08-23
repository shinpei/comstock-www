package model

import (
	"time"
)

type User struct {
	Mail      string
	Username  string
	UID       int
	Created   int64 // TODO: make it with time package
	LastLogin int64 // TODO: make it with time package
}

func CreateUser(mail string, username string, uid int) *User {
	return &User{Mail: mail, Username: username, UID: uid}
}
func CreateUserForNewCommer(mail string, uid int) *User {
	return &User{Mail: mail, UID: uid, Created: time.Now().Unix(), LastLogin: time.Now().Unix()}
}

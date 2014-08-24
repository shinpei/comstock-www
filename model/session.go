package model

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Session struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Token   string
	UID     int
	Expires int64 // TODO: replace it with time
}

func CreateSession(token string, uid int) *Session {
	return &Session{Token: token, UID: uid, Expires: time.Now().Unix()}
}
func UpdateSessionToken(id bson.ObjectId, token string, uid int) *Session {
	return &Session{ID: id, Token: token, UID: uid, Expires: time.Now().Unix()}
}

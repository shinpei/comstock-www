package model

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Session struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Token   string
	UID     int
	Expires time.Time // TODO: replace it with time
}

func CreateNewSession(token string, uid int) *Session {
	return &Session{ID: bson.NewObjectId(), Token: token, UID: uid, Expires: time.Now()}
}
func UpdateSessionToken(id *bson.ObjectId, token string, uid int) *Session {
	return &Session{ID: *id, Token: token, UID: uid, Expires: time.Now()}
}

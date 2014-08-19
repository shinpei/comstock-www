package model

type Session struct {
	Token   string
	UID     int
	Expires int // TODO: replace it with time
}

func CreateSession(token string, uid int) *Session {
	return &Session{Token: token, UID: uid}
}

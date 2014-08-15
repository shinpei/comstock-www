package model

type LoginInfo struct {
	mail string
	pass string
}

func CreateLoginInfo(m string, p string) *LoginInfo {
	return &LoginInfo{mail: m, pass: p}
}

func (l *LoginInfo) Mail() string {
	return l.mail
}

func (l *LoginInfo) Pass() string {
	return l.pass
}

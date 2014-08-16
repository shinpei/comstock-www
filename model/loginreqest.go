package model

type LoginRequest struct {
	mail string
	pass string
}

func CreateLoginRequest(m string, p string) *LoginRequest {
	return &LoginRequest{mail: m, pass: p}
}

func (l *LoginRequest) Mail() string {
	return l.mail
}

func (l *LoginRequest) Pass() string {
	return l.pass
}

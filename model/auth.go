package model

type Auth struct {
	UID      int
	Date     string // TODO: make it with time pkg
	Password string
}

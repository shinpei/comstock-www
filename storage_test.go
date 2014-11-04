package main

import (
	"github.com/shinpei/comstock-www/model"
	"testing"
	"time"
)

func TestPostHistory(t *testing.T) {

	ses, db := getSessionAndDB()
	defer ses.Close()
	s, err := loginAs(db, model.CreateLoginRequest("test@mail.com", "test"))

	if err != nil {
		println("error", err.Error())
		panic("cannot login test db")
	}
	accessToken := s.Token

	err = postHistory(db, accessToken, "ls -la", time.Now(), "hi")
	if err != nil {
		println("error", err.Error())
	}
}

func TestFetchHistory(t *testing.T) {
	ses, db := getSessionAndDB()
	defer ses.Close()
	s, err := loginAs(db, model.CreateLoginRequest("test@mail.com", "test"))
	if err != nil {
		println(err.Error())
		panic("cannot login test db")
	}
	accessToken := s.Token
	_ = accessToken

}

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
	tk := s.Token

	err = postHistory(db, tk, "ls -la", time.Now(), "hi")
	if err != nil {
		println("error", err.Error())
	}
}

func TestPostHistories(t *testing.T) {
	ses, db := getSessionAndDB()
	defer ses.Close()
	s, err := loginAs(db, model.CreateLoginRequest("test@mail.com", "test"))
	if err != nil {
		println(err.Error())
		panic("cannot login test db")
	}
	tk := s.Token
	_ = tk
	_, ci1 := model.CreateNewCommandItem("ls -la")
	_, ci2 := model.CreateNewCommandItem("wc -l")
	cis := []*model.NewCommandItem{ci1, ci2}
	_, f := model.CreateFlow(cis)
	h := model.CreateHistoryFromFlow(1, time.Now(), "sample", f)
	err = InsertHistory(db, h)
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
	tk := s.Token
	hist, err := FindHistoryFromNum(db, tk, 2)
	if err != nil || hist == nil {
		println(err.Error())
		t.Fail()
	}
	println(hist.Command())
}

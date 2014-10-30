// test for model translation from 1 to 2
package model

import (
	. "github.com/shinpei/comstock/test"
	"testing"
	"time"
)

func createCommandList() []*NewCommandItem {
	_, ci1 := CreateNewCommandItem("ls -la")
	_, ci2 := CreateNewCommandItem("wc -l")
	return []*NewCommandItem{ci1, ci2}
}

func TestCreateNewCommandItem(t *testing.T) {
	_, ci := CreateNewCommandItem("ls -la")
	AssertEqual(t, "ls -la", ci.Command)
}

func TestCreateFlow(t *testing.T) {
	_, f := CreateFlow(createCommandList())
	AssertEqual(t, "ls -la => wc -l", f.Command())
}

func TestCreateHistoryFromFlow(t *testing.T) {
	_, f := CreateFlow(createCommandList())
	h := CreateHistoryFromFlow(1, time.Now(), "sample", f)
	AssertEqual(t, "ls -la => wc -l", h.Command())
}

func TestTrans(t *testing.T) {
	c1 := CreateCommandItem(1, "ls -la")
	history := TranslateCommand1to2(c1)
	AssertEqual(t, c1.Data.Command, history.Command())
	AssertEqual(t, c1.UID, history.UID)
}

func TestHistoryCommand(t *testing.T) {
	c1 := CreateCommandItem(1, "ls -la")
	history := TranslateCommand1to2(c1)
	AssertEqual(t, "ls -la", history.Command())
}

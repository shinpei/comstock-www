// test for model translation from 1 to 2
package model

import (
	. "github.com/shinpei/comstock/test"
	"testing"
	"time"
)

func createCommandList() []*CommandItem {
	ci1 := CreateCommandItem("ls -la")
	ci2 := CreateCommandItem("wc -l")
	return []*CommandItem{ci1, ci2}
}

func TestCreateCommandItem(t *testing.T) {
	ci := CreateCommandItem("ls -la")
	AssertEqual(t, "ls -la", ci.Command)
}

func TestCreateFlow(t *testing.T) {
	f := CreateFlow(createCommandList())
	AssertEqual(t, "ls -la => wc -l", f.Command())
}

func TestCreateHistoryFromFlow(t *testing.T) {
	f := CreateFlow(createCommandList())
	h := CreateHistoryFromFlow(1, time.Now(), "sample", f)
	AssertEqual(t, "ls -la => wc -l", h.Command())
}

func TestTrans(t *testing.T) {
	c1 := CreateOldCommandItem(1, "ls -la")
	history := TranslateCommand1to2(c1)
	AssertEqual(t, c1.Data.Command, history.Command())
	AssertEqual(t, c1.UID, history.UID)
}

func TestHistoryCommand(t *testing.T) {
	c1 := CreateOldCommandItem(1, "ls -la")
	history := TranslateCommand1to2(c1)
	AssertEqual(t, "ls -la", history.Command())
}

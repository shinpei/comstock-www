// test for model translation from 1 to 2
package model

import (
	"fmt"
	. "github.com/shinpei/comstock/test"
	"testing"
)

func TestCreateNewCommandItem(t *testing.T) {
	_, ci := CreateNewCommandItem(1, "ls -la")
	AssertEqual(t, "ls -la", ci.Command)
}

func TestCreateFlow(t *testing.T) {
	_, ci1 := CreateNewCommandItem(1, "ls -la")
	_, ci2 := CreateNewCommandItem(1, "wc -l")
	_, f := CreateFlow([]*NewCommandItem{ci1, ci2})

	fmt.Println(f.Command())
}

func TestTrans(t *testing.T) {
	c1 := CreateCommandItem(1, "ls -la")
	history := TranslateCommand1to2(c1)
	AssertEqual(t, c1.Data.Command, history.FlowPtr.ItemsPtr[0].Command)
	AssertEqual(t, c1.UID, history.UID)
	fmt.Printf("%#v, %#v\n", c1.Date, history.Date)
}

func TestHistoryCommand(t *testing.T) {
	c1 := CreateCommandItem(1, "ls -la")
	history := TranslateCommand1to2(c1)
	fmt.Printf("cmd:%s\n", history.Command())

}

// test for model translation from 1 to 2
package model

import (
	"fmt"
	. "github.com/shinpei/comstock/test"
	"testing"
)

func TestTrans(t *testing.T) {
	c1 := CreateCommandItem(1, "ls -la")
	history := TranslateCommand1to2(c1)
	AssertEqual(t, c1.Data.Command, history.FlowPtr.ItemsPtr[0].Command)
	AssertEqual(t, c1.UID, history.UID)
	fmt.Printf("%#v, %#v\n", c1.Date, history.Date)
}

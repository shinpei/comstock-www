// test for model translation from 1 to 2
package model

import (
	"fmt"
	. "github.com/shinpei/comstock/test"
	"testing"
)

func TestTrans(t *testing.T) {
	c1 := CreateCommandItem(1, "ls -la")
	c2 := TranslateCommand1to2(c)
	AssertEqual(t, c1.Command, c2.Command)
	AssertEqual(t, c1.UID, c2.UID)
	fmt.Printf("%#v, %#v\n", c.Date, c2.Date)
}

// test for model translation from 1 to 2
package model

import (
	"fmt"
	. "github.com/shinpei/comstock/test"
	"testing"
)

func TestTrans(t *testing.T) {
	c := CreateCommandItem(1, "ls -la")
	c2 := TranslateCommand1to2(c)
	AssertEqual(t, "ls -la", c2.Command)
	fmt.Printf("%#v, %#v\n", c.Date, c2.Date)

}

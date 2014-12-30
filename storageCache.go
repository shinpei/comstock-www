package main

import (
	"sync"
)

var fcount int = 0
var flowPool = sync.Pool{
	New: func() interface{} {
		fcount++
		return &flow{}
	},
}
var hcount int = 0
var histPool = sync.Pool{
	New: func() interface{} {
		hcount++
		return &history{}
	},
}

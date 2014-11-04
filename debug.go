// +build debug

package main

func D(fmt string, argv ...interface{}) {
	println("[DEBUG]"+fmt, argv)
}

// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "github.com/beta/tigo"

func main() {
	win := tigo.NewWindow(320, 240, "Hello TIGO", tigo.WindowAuto)
	for !win.Closed() {
		win.Clear(tigo.RGB(0x80, 0x90, 0xa0))
		win.Print(tigo.DefaultFont(), 120, 110, tigo.RGB(0xff, 0xff, 0xff), "Hello TIGO")
		win.Update()
	}
	win.Free()
}

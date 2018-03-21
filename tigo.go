// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package tigo implements the interface for TIGO.
package tigo

// #cgo darwin LDFLAGS: -framework OpenGL -framework Cocoa
// #include "tigr/tigr.h"
// #include "tigr/tigr.c"
import "C"
import "unsafe"

// Pixel represents one pixel.
type Pixel struct {
	B uint8
	G uint8
	R uint8
	A uint8
}

type tPixel _Ctype_TPixel

// WindowFlag reprensets a flag for window sizing mode.
type WindowFlag int

const (
	// WindowFixed means a window's bitmap has a fixed size.
	WindowFixed WindowFlag = 1 << iota
	// WindowAuto means a window's bitmap will automatically resize after each Update.
	WindowAuto
	// Window2X enforces (at least) 2X pixel scale.
	Window2X
	// Window3X enforces (at least) 3X pixel scale.
	Window3X
	// Window4X enforces (at least) 4X pixel scale.
	Window4X
	// WindowRetina enables retina support on macOS.
	WindowRetina
)

// Bitmap represents a bitmap.
type Bitmap struct {
	cBitmap unsafe.Pointer
}

// NewWindow creates a new empty window. title is UTF-8.
func NewWindow(width, height int, title string, flags WindowFlag) *Bitmap {
	tigr := C.tigrWindow(C.int(width), C.int(height), C.CString(title), C.int(flags))
	return &Bitmap{
		cBitmap: unsafe.Pointer(tigr),
	}
}

// NewBitmap creates an empty off-screen bitmap.
func NewBitmap(width, height int) *Bitmap {
	tigr := C.tigrBitmap(C.int(width), C.int(height))
	return &Bitmap{
		cBitmap: unsafe.Pointer(tigr),
	}
}

// Free deletes a window/bitmap.
func (bmp *Bitmap) Free() {
	C.tigrFree((*C.Tigr)(bmp.cBitmap))
}

// Closed returns true if the user requested to close a window.
func (bmp *Bitmap) Closed() bool {
	closed := int(C.tigrClosed((*C.Tigr)(bmp.cBitmap)))
	return closed > 0
}

// Update displays a window's content onto screen.
func (bmp *Bitmap) Update() {
	C.tigrUpdate((*C.Tigr)(bmp.cBitmap))
}

// SetPostFX sets post-FX properties for a window.
// hBlur/vBlur = whether to use bilinear filtering along that axis.
// scanlines = CRT scanlines effect (0-1).
// contrast = contrast boost (1 = no change, 2 = 2X contrast, etc)
func (bmp *Bitmap) SetPostFX(hBlur, vBlur bool, scanlines, contrast float32) {
	var hBlurInt, vBlurInt int
	if hBlur {
		hBlurInt = 1
	}
	if vBlur {
		vBlurInt = 1
	}

	C.tigrSetPostFX((*C.Tigr)(bmp.cBitmap), C.int(hBlurInt), C.int(vBlurInt), C.float(scanlines), C.float(contrast))
}

// Drawing

// Get gets a pixel from bitmap.
func (bmp *Bitmap) Get(x, y int) Pixel {
	cPixel := C.tigrGet((*C.Tigr)(bmp.cBitmap), C.int(x), C.int(y))
	return Pixel{
		B: uint8(cPixel.b),
		G: uint8(cPixel.g),
		R: uint8(cPixel.r),
		A: uint8(cPixel.a),
	}
}

// Plot sets a pixel onto bitmap.
func (bmp *Bitmap) Plot(x, y int, p Pixel) {
	cPixel := tPixel{
		b: C.uchar(p.B),
		g: C.uchar(p.G),
		r: C.uchar(p.R),
		a: C.uchar(p.A),
	}
	C.tigrPlot((*C.Tigr)(bmp.cBitmap), C.int(x), C.int(y), C.TPixel(cPixel))
}

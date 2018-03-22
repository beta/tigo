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

// goPixel converts a TPixel from C into a Pixel from Go.
func goPixel(p C.TPixel) Pixel {
	return Pixel{
		B: uint8(p.b),
		G: uint8(p.g),
		R: uint8(p.r),
		A: uint8(p.a),
	}
}

// tPixel converts a Pixel from Go into a TPixel from C.
func cPixel(p Pixel) C.TPixel {
	return C.TPixel{
		b: C.uchar(p.B),
		g: C.uchar(p.G),
		r: C.uchar(p.R),
		a: C.uchar(p.A),
	}
}

// RGB is a helper function for making colors.
func RGB(r, g, b uint8) Pixel {
	return Pixel{
		R: r,
		G: g,
		B: b,
		A: 0xff,
	}
}

// RGBA is a helper function for making colors with alpha value.
func RGBA(r, g, b, a uint8) Pixel {
	return Pixel{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

// Get gets a pixel from bitmap.
func (bmp *Bitmap) Get(x, y int) Pixel {
	return goPixel(C.tigrGet((*C.Tigr)(bmp.cBitmap), C.int(x), C.int(y)))
}

// Plot sets a pixel onto bitmap.
func (bmp *Bitmap) Plot(x, y int, p Pixel) {
	C.tigrPlot((*C.Tigr)(bmp.cBitmap), C.int(x), C.int(y), cPixel(p))
}

// Clear clears a bitmap to a color.
func (bmp *Bitmap) Clear(color Pixel) {
	C.tigrClear((*C.Tigr)(bmp.cBitmap), cPixel(color))
}

// Fill fills in a solid rectangle.
func (bmp *Bitmap) Fill(x, y, w, h int, color Pixel) {
	C.tigrFill((*C.Tigr)(bmp.cBitmap), C.int(x), C.int(y), C.int(w), C.int(h), cPixel(color))
}

// Rect draws an empty rectangle (exclusive coords).
func (bmp *Bitmap) Rect(x, y, w, h int, color Pixel) {
	C.tigrRect((*C.Tigr)(bmp.cBitmap), C.int(x), C.int(y), C.int(w), C.int(h), cPixel(color))
}

// Line draws a line.
func (bmp *Bitmap) Line(x0, y0, x1, y1 int, color Pixel) {
	C.tigrLine((*C.Tigr)(bmp.cBitmap), C.int(x0), C.int(y0), C.int(x1), C.int(y1), cPixel(color))
}

// Blit copies bitmap data from bmp to dest.
// sx/sy = source coordinates.
// dx/dy = dest coordinates.
// w/h: width/height.
func (bmp *Bitmap) Blit(sx, sy int, dest *Bitmap, dx, dy, w, h int) {
	C.tigrBlit((*C.Tigr)(dest.cBitmap), (*C.Tigr)(bmp.cBitmap), C.int(dx), C.int(dy), C.int(sx), C.int(sy), C.int(w), C.int(h))
}

// BlitAlpha is same as Blit, but blends with the bitmap alpha channel, and uses the 'alpha' variable to fade out.
func (bmp *Bitmap) BlitAlpha(sx, sy int, dest *Bitmap, dx, dy, w, h int, alpha float32) {
	C.tigrBlitAlpha((*C.Tigr)(dest.cBitmap), (*C.Tigr)(bmp.cBitmap), C.int(dx), C.int(dy), C.int(sx), C.int(sy), C.int(w), C.int(h), C.float(alpha))
}

// BlitTint is same as Blit, but tints the source bitmap with a color.
func (bmp *Bitmap) BlitTint(sx, sy int, dest *Bitmap, dx, dy, w, h int, tint Pixel) {
	C.tigrBlitTint((*C.Tigr)(dest.cBitmap), (*C.Tigr)(bmp.cBitmap), C.int(dx), C.int(dy), C.int(sx), C.int(sy), C.int(w), C.int(h), cPixel(tint))
}

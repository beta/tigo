// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package tigo implements the interface for TIGO.
package tigo

// #cgo darwin LDFLAGS: -framework OpenGL -framework Cocoa
// #cgo windows LDFLAGS: -ld3d9
// #include "tigr/tigr.h"
// #include "tigr/tigr.c"
//
// void wrapTigrText(Tigr *dest, TigrFont *font, int x, int y, TPixel color, const char *text) {
//     tigrPrint(dest, font, x, y, color, text);
// }
import "C"
import (
	"fmt"
	"unsafe"
)

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
	return &Bitmap{unsafe.Pointer(tigr)}
}

// NewBitmap creates an empty off-screen bitmap.
func NewBitmap(width, height int) *Bitmap {
	tigr := C.tigrBitmap(C.int(width), C.int(height))
	return &Bitmap{unsafe.Pointer(tigr)}
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

// Font printing

// Glyph represents a glyph.
type Glyph struct {
	Code int
	X    int
	Y    int
	W    int
	H    int
}

// Font represents a font consisting of glyphs.
type Font struct {
	cFont unsafe.Pointer
}

// Codepage is a number representing a codepage.
type Codepage int

// Codepages.
const (
	// ASCII represents regular 7-bit ASCII codepage.
	ASCII Codepage = 0
	// Windows1252 represents Windows 1252 codepage.
	Windows1252 Codepage = 1252
)

// LoadFont loads a font.
// The font bitmap should contain all characters for the given codepage, excluding the first 32 control codes.
// Supported codepages:
//     0    - Regular 7-bit ASCII
//     1252 - Windows 1252
func LoadFont(bmp *Bitmap, codepage Codepage) *Font {
	cFont := C.tigrLoadFont((*C.Tigr)(bmp.cBitmap), C.int(codepage))
	return &Font{unsafe.Pointer(cFont)}
}

// Print prints UTF-8 text onto a bitmap.
func (bmp *Bitmap) Print(font *Font, x, y int, color Pixel, text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.wrapTigrText((*C.Tigr)(bmp.cBitmap), (*C.TigrFont)(font.cFont), C.int(x), C.int(y), cPixel(color), cText)
}

// TextWidth returns the width of a string.
func (font *Font) TextWidth(text string) int {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	return int(C.tigrTextWidth((*C.TigrFont)(font.cFont), cText))
}

// TextHeight returns the height of a string.
func (font *Font) TextHeight(text string) int {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	return int(C.tigrTextHeight((*C.TigrFont)(font.cFont), cText))
}

// DefaultFont returns the default built-in font.
func DefaultFont() *Font {
	return &Font{unsafe.Pointer(C.tfont)}
}

// User input.

// Key represents a key scancode. ASCII ('A'-'Z' and '0'-'9') is used for letters/numbers.
type Key int

// Key scancode constants.
const (
	Pad0 Key = 128 + iota
	Pad1
	Pad2
	Pad3
	Pad4
	Pad5
	Pad6
	Pad7
	Pad8
	Pad9
	PadMul
	PadAdd
	PadEnter
	PadSub
	PadDot
	PadDiv
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	Backspace
	Tab
	Return
	Shift
	Control
	Alt
	Pause
	CapsLock
	Escape
	Space
	PageUp
	PageDown
	End
	Home
	Left
	Up
	Right
	Down
	Insert
	Delete
	LWin
	RWin
	NumLock
	ScrollLock
	LShift
	RShift
	LControl
	RControl
	LAlt
	RAlt
	SemiColon
	Equals
	Comma
	Minus
	Dot
	Slash
	BackTick
	LSquare
	BackSlash
	RSquare
	Tick
)

// Mouse returns mouse input for a window.
func (bmp *Bitmap) Mouse() (x, y, buttons int) {
	var cX, cY, cButtons C.int
	C.tigrMouse((*C.Tigr)(bmp.cBitmap), (*C.int)(&cX), (*C.int)(&cY), (*C.int)(&cButtons))
	x, y, buttons = int(cX), int(cY), int(cButtons)
	return
}

// KeyDown returns true if a key is pressed for a window.
// KeyDown only tests for the initial press.
func (bmp *Bitmap) KeyDown(key Key) bool {
	return int(C.tigrKeyDown((*C.Tigr)(bmp.cBitmap), C.int(key))) > 0
}

// KeyHeld returns true if a key is held for a window.
// KeyHeld repreats each frame.
func (bmp *Bitmap) KeyHeld(key Key) bool {
	return int(C.tigrKeyHeld((*C.Tigr)(bmp.cBitmap), C.int(key))) > 0
}

// ReadChar reads character input for a window and returns the Unicode value of the last key pressed.
// If no key is pressed, ReadChar returns 0.
func (bmp *Bitmap) ReadChar() int {
	return int(C.tigrReadChar((*C.Tigr)(bmp.cBitmap)))
}

// Bitmap I/O.

// LoadImage loads a PNG, from either a file. fileName is UTF-8.
func LoadImage(fileName string) (*Bitmap, error) {
	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	tigr := C.tigrLoadImage(cFileName)
	if tigr == nil {
		return nil, fmt.Errorf("failed to load image")
	}
	return &Bitmap{unsafe.Pointer(tigr)}, nil
}

// SaveImage saves a PNG to a file. fileName is UTF-8.
func SaveImage(fileName string, bmp *Bitmap) error {
	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	result := C.tigrSaveImage(cFileName, (*C.Tigr)(bmp.cBitmap))
	if int(result) == 0 {
		return fmt.Errorf("failed to save image")
	}
	return nil
}

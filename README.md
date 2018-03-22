# TIGO

TIGO is a Go binding for [TIGR](https://bitbucket.org/rmitton/tigr), a tiny graphics library.

TIGO is still working in progress, and has been tested only on macOS. More features and better support for other platforms will be comming soon.

## Installation

```bash
$ go get github.com/beta/tigo
```

## Example

Below is a basic example of creating window and drawing text with TIGO.

```go
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
```

More examples can be found in the *example* directory.

## License

MIT

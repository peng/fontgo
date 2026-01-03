# Font Go

Font Go is a Go library for reading, subsetting, and writing font files. It supports various font formats TTF.

```go
package main

import (
    "fmt"
    "github.com/peng/fontgo/font"
)

func main() {
    f := font.ReadFontFile('font.ttf')
    defer func() {
        if err := f.Close(); err != nil {
            fmt.Println(err)
        }
    }

    // will get font struct
    f.GetFontInfo()

    // only get some word font
    f.Subset(["a", "b"])

    // write font
    if err := f.Write("font.woff"); err != nil {
        fmt.Println(err)
    }
}

```
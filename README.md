# Font Go

```go
package main

import (
    "fmt"
    "github.com/peng/fontgo"
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

    // convert font format
    f.Convert("svg", opt)

    // write font
    if err := f.Write("font.woff"); err != nil {
        fmt.Println(err)
    }
}

```
package main

import (
	"fmt"

	"github.com/peng/fontgo/font"
)

func main() {
	dir, err := font.DataReader("./test/HanyiSentyCrayon.ttf")
	if (err == nil) {
		fmt.Printf("%v",dir)
	}
}

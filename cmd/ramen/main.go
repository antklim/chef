package main

import (
	"fmt"

	"github.com/antklim/ramen"
)

func main() {
	fmt.Println("Ramen v0.1.0")
	p := ramen.New("XYZ")
	p.Init("")
}

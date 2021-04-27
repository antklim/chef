package main

import (
	"fmt"

	"github.com/antklim/ramen"
)

// TODO: Add CLI options support
// TODO: Wire with ramen to init project

func main() {
	fmt.Println("Ramen v0.1.0")
	p := ramen.New("XYZ")
	p.Init("")
}

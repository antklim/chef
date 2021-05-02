package main

import (
	"fmt"

	"github.com/antklim/chef"
)

// TODO: Add CLI options support
// TODO: Wire with chef to init project

func main() {
	fmt.Println("Chef v0.1.0")
	p := chef.New("XYZ")
	p.Init("")
}

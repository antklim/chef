package main

import (
	"fmt"
	"log"

	"github.com/antklim/chef"
)

// TODO: Add CLI options support
// TODO: Wire with chef to init project

func main() {
	fmt.Println("Chef v0.1.0")
	p := chef.New("XYZ")
	if err := p.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := p.Init(); err != nil {
		log.Fatal(err)
	}
}
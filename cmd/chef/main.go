package main

import (
	"fmt"
	"log"

	"github.com/antklim/chef"
)

// TODO: Add CLI options support

func main() {
	fmt.Println("Chef v0.1.0")
	if err := chef.Init("XYZ"); err != nil {
		log.Fatal(err)
	}
}

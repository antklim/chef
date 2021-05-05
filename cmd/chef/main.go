package main

import (
	"fmt"
	"os"

	"github.com/antklim/chef"
)

// TODO: Add CLI options support

func main() {
	fmt.Println("Chef v0.1.0")
	if err := chef.Init("XYZ"); err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}

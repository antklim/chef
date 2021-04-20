package ramen

import (
	"os"
)

// Project manager.
type Project struct {
}

// New project.
func New() *Project {
	return &Project{}
}

// Init initializes the project layout.
func (p *Project) Init(name string) error {
	return os.Mkdir(name, 0755)
}

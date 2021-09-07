package cli

import (
	"io"
	"os"
)

var printout io.Writer = os.Stdout

type Project interface {
	Init() error
	Build() (string, error)
	ComponentsNames() []string // TODO: replace with components
	EmployComponent(string, string) error
}

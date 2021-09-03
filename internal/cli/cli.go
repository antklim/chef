package cli

import (
	"io"
	"os"
)

var printout io.Writer = os.Stdout

type Project interface {
	Init() error
	Build() (string, error)
	ComponentsNames() []string
	EmployComponent(string, string) error
}

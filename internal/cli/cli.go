package cli

import (
	"io"
	"os"

	"github.com/antklim/chef/internal/project"
)

var printout io.Writer = os.Stdout

type Project interface {
	Init() error
	Build() (string, error)
	Components() []project.Component
	EmployComponent(string, string) error
}

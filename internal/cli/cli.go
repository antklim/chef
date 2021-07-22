package cli

import "github.com/antklim/chef/internal/project"

type Project interface {
	Bootstrap() error
	Add(project.Component) error
	Location() (string, error)
	Name() string
}

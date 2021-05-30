package chef

import (
	"github.com/antklim/chef/internal/project"
)

// TODO: delete it after init moved to cli.bootstrap

// Init initializes default project layout.
func Init(name string) error {
	p := project.New(name)
	if err := p.Validate(); err != nil {
		return err
	}
	return p.Init()
}

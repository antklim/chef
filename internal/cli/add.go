package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO: implement add project component command
func addCmd() *cobra.Command { // nolint
	return nil
}

func AddCmdRunner(p Project, component, name string) error {
	if err := p.Add(component, name); err != nil {
		return errors.Wrap(err, "unable to add to a project")
	}
	return nil
}
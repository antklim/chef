package cli

import (
	"os"
	"path"

	"github.com/antklim/chef/internal/chef"
	"github.com/antklim/chef/internal/display"
	"github.com/antklim/chef/internal/project"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO: by default employ component adds a component to the current
// directory (assume current directory is a root of the project).
// Add ability to provide a project name and location.

// TODO: add register component

var (
	component = Flag{
		LongForm:   "component",
		ShortForm:  "c",
		Help:       "The component to employ.",
		IsRequired: true,
	}
	componentName = Flag{
		LongForm:   "name",
		ShortForm:  "n",
		Help:       "Name of the node to be created employing the component.",
		IsRequired: true,
	}
)

func componentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "components",
		Short: "Manage project components",
		Long:  "Manage project components",
	}

	cmd.AddCommand(listComponentsCmd())
	cmd.AddCommand(employComponentCmd())

	return cmd
}

func listComponentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Short:   "List project components",
		Long:    "List registered project components",
		Example: `chef components list
chef components ls`,
		RunE: func(_ *cobra.Command, _ []string) error {
			p, err := initProject()
			if err != nil {
				return err
			}
			return componentsListCmdRunner(p)
		},
	}

	return cmd
}

func employComponentCmd() *cobra.Command {
	var inputs struct {
		Component string // component name
		Name      string // node name to be created using the component
	}

	cmd := &cobra.Command{
		Use:   "employ",
		Args:  cobra.NoArgs,
		Short: "Employ project component",
		Long:  "Use component to add a new functionality to a project",
		Example: `chef components employ --component http_handler --name foo 
chef components employ -c http_handler -n bar`,
		RunE: func(_ *cobra.Command, _ []string) error {
			p, err := initProject()
			if err != nil {
				return err
			}
			return componentsEmployCmdRunner(p, inputs.Component, inputs.Name)
		},
	}

	component.RegisterString(cmd, &inputs.Component, "")
	componentName.RegisterString(cmd, &inputs.Name, "")

	return cmd
}

func componentsListCmdRunner(p Project) error {
	if err := p.Init(); err != nil {
		return errors.Wrap(err, "init project failed")
	}

	return display.ComponentsList(printout, p.Components())
}

func componentsEmployCmdRunner(p Project, component, name string) error {
	if err := p.Init(); err != nil {
		return errors.Wrap(err, "init project failed")
	}

	if err := p.EmployComponent(component, name); err != nil {
		// TODO: better explanation why employ failed
		return errors.Wrapf(err, "employ %q component failed", component)
	}

	return display.ComponentsEmploy(printout, name, component)
}

func initProject() (*project.Project, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get working directory")
	}

	f, err := os.Open(path.Join(dir, chef.DefaultNotationFileName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open notation")
	}

	n, err := chef.ReadNotation(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read notation")
	}

	p := project.New(path.Base(dir), project.WithRoot(path.Dir(dir)), project.WithNotation(n))
	return p, nil
}

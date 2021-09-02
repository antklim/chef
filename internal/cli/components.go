package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/antklim/chef/internal/project"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO (feat): add register component
// TODO (ref): add command tests for employ component

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
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO (feat): implement list registered components
			fmt.Println("not implemented")
			return nil
		},
	}

	return cmd
}

// TODO (feat): by default this command add component to a current directory (assume current directory is a root of the project)
// add ability to provide a project name and location

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
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to get working directory")
			}
			p := project.New(path.Base(dir), project.WithRoot(path.Dir(dir)), project.WithServer("http"))
			return componentsEmployCmdRunner(p, inputs.Component, inputs.Name)
		},
	}

	component.RegisterString(cmd, &inputs.Component, "")
	componentName.RegisterString(cmd, &inputs.Name, "")

	return cmd
}

func componentsEmployCmdRunner(p Project, component, name string) error {
	if err := p.Init(); err != nil {
		return errors.Wrap(err, "init project failed")
	}

	if err := p.EmployComponent(component, name); err != nil {
		// TODO: better explanation why employ failed
		return errors.Wrapf(err, "employ %q component failed", component)
	}

	fmt.Printf("successfully added %q as %q component\n", name, component)

	return nil
}

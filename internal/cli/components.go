package cli

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

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
			// TODO: implement list registered components
			fmt.Println("not implemented")
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement employ component
			fmt.Println("not implemented")
			return nil
		},
	}

	component.RegisterString(cmd, &inputs.Component, "")
	componentName.RegisterString(cmd, &inputs.Name, "")

	return cmd
}

func componentsEmployCmdRunner(p Project, component, name string) error {
	if err := p.Employ(component, name); err != nil {
		return errors.Wrap(err, "project employ component failed")
	}
	return nil
}

package cli

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	// var inputs struct {
	// 	Component string // component name
	// 	Name      string // project layout node name to be created using the component
	// }

	cmd := &cobra.Command{
		Use:   "employ",
		Args:  cobra.NoArgs,
		Short: "Employ project component",
		Long:  "Use component to add a new functionality to a project",
		Example: `chef components employ --component http_handler --name foo 
chef components employ -c http_handler -n bar`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement list registered components
			fmt.Println("not implemented")
			return nil
		},
	}

	return cmd
}

func componentsEmployCmdRunner(p Project, component, name string) error {
	if err := p.Add(component, name); err != nil {
		return errors.Wrap(err, "could not add project component")
	}
	return nil
}

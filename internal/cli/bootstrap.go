package cli

import (
	"fmt"

	"github.com/antklim/chef"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// var projectCategoryOptions = []string{
// 	"CLI",
// 	"Package",
// 	"Service",
// }

func bootstrapCmd() *cobra.Command {
	// Options
	// Name: project name
	// Root: project location root (where to create a project)
	// Category: pkg, app, cli
	// Layout: TBA

	var inputs struct {
		Name     string
		Root     string
		Category string
		Layout   string
	}

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Args:  cobra.NoArgs,
		Short: "Bootstrap a new project",
		Long:  "Bootstrap a new project",
		Example: `chef bootstrap --name myproject
chef bootstrap --category [cli|pkg|srv] --name myproject
chef bootstrap -c [cli|pkg|srv] -n myproject --root /usr/local --layout chef.yml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// project := project {
			// 	Name: inputs.Name,
			// 	...
			// }
			fmt.Println("Chef v0.1.0")
			return chef.Init("XYZ")
		},
	}

	cmd.Flags().StringVarP(&inputs.Name, "name", "n", "", "Name of the project.")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(errors.Wrap(err, "failed to register string flag"))
	}

	cmd.Flags().StringVarP(&inputs.Root, "root", "r", "", "Root location of the project.")

	cmd.Flags().StringVarP(&inputs.Category, "category", "c", "", "Category of project:\n"+
		"- cli: CLI application.\n"+
		"- pkg: package.\n"+
		"- srv: service application based on HTTP or gRPC.\n")

	cmd.Flags().StringVarP(&inputs.Layout, "layout", "l", "", "Location of the project layout configuration.")

	return cmd
}

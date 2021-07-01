package cli

import (
	"fmt"

	"github.com/antklim/chef/internal/project"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO: add commamds tests

var (
	projName = Flag{
		LongForm:   "name",
		ShortForm:  "n",
		Help:       "Name of the project.",
		IsRequired: true,
	}
	projRoot = Flag{
		LongForm:   "root",
		ShortForm:  "r",
		Help:       "Root location of the project.",
		IsRequired: false,
	}
	projCategory = Flag{
		LongForm:  "category",
		ShortForm: "c",
		Help: "Category of project:\n" +
			"- cli: CLI application.\n" +
			"- pkg: package.\n" +
			"- srv: service application based on HTTP or gRPC.\n",
		IsRequired: true,
	}
	projModule = Flag{
		LongForm:   "module",
		ShortForm:  "m",
		Help:       "Name of the project's module to be used in 'go mod'.",
		IsRequired: true,
	}
	projLayout = Flag{
		LongForm:   "layout",
		ShortForm:  "l",
		Help:       "Location of the project's layout configuration.",
		IsRequired: false,
	}
	projServer = Flag{
		LongForm:   "server",
		ShortForm:  "s",
		Help:       "Server type for projects of category service.",
		IsRequired: false,
	}
)

func bootstrapCmd() *cobra.Command {
	var inputs struct {
		Name     string
		Root     string
		Category string
		Module   string
		Layout   string
		Server   string
	}

	cmd := &cobra.Command{
		Use:   "boot",
		Args:  cobra.NoArgs,
		Short: "Bootstrap a new project",
		Long:  "Bootstrap a new project",
		Example: `chef boot --name myproject
chef boot --category [cli|pkg|srv] --name myproject
chef boot -c [cli|pkg|srv] -n myproject --root /usr/local --layout chef.yml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := project.New(inputs.Name,
				project.WithRoot(inputs.Root),
				project.WithCategory(inputs.Category),
				project.WithServer(inputs.Server),
				project.WithModule(inputs.Module),
				// TODO: layout location
			)

			if err := p.Bootstrap(); err != nil {
				return errors.Wrap(err, "unable to bootstrap project")
			}

			fmt.Printf("project %s successfully bootrapped\n", p.Name())

			if l, err := p.Location(); err != nil {
				fmt.Printf("unable to get project location: %+v\n", err)
			} else {
				fmt.Printf("project location: %s\n", l)
			}

			return nil
		},
	}

	projName.RegisterString(cmd, &inputs.Name, "")
	projRoot.RegisterString(cmd, &inputs.Root, "")
	projCategory.RegisterString(cmd, &inputs.Category, "")
	projModule.RegisterString(cmd, &inputs.Module, "")
	projLayout.RegisterString(cmd, &inputs.Layout, "")
	projServer.RegisterString(cmd, &inputs.Server, "")

	return cmd
}

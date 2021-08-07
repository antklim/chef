package cli

import (
	"fmt"

	"github.com/antklim/chef/internal/project"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO: add commamds tests, use https://github.com/commander-cli/commander
// TODO: layout/chef.yml location

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

func initCmd() *cobra.Command {
	var inputs struct {
		Name     string
		Root     string
		Category string
		Module   string
		Layout   string
		Server   string
	}

	cmd := &cobra.Command{
		Use:   "init",
		Args:  cobra.NoArgs,
		Short: "Initialize a new project",
		Long:  "initialize a new project",
		Example: `chef init --name myproject
chef init --category [srv] --name myproject
chef init -c [srv] -n myproject --root /usr/local`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := project.New(inputs.Name,
				project.WithRoot(inputs.Root),
				project.WithCategory(inputs.Category),
				project.WithServer(inputs.Server),
				project.WithModule(inputs.Module),
			)
			return initCmdRunner(p)
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

func initCmdRunner(p Project) error {
	if err := p.Init(); err != nil {
		// TODO: don't print the stack trace
		return errors.Wrap(err, "init project failed")
	}

	loc, err := p.Build()
	if err != nil {
		return errors.Wrap(err, "init project failed")
	}

	fmt.Printf("project successfully inited at %s\n", loc)

	return nil
}

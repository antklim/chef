package cli

import (
	"fmt"

	"github.com/antklim/chef/internal/project"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

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
	projLayout = Flag{
		LongForm:   "layout",
		ShortForm:  "l",
		Help:       "Location of the project layout configuration.",
		IsRequired: false,
	}
)

func bootstrapCmd() *cobra.Command {
	var inputs struct {
		Name     string
		Root     string
		Category string
		Layout   string
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
			projCategory := project.CategoryFor(inputs.Category)

			if projCategory.IsUnknown() {
				return fmt.Errorf("unknown project category: %s", inputs.Category)
			}

			p := project.New(inputs.Name,
				project.WithRoot(inputs.Root),
				project.WithCategory(project.Category(inputs.Category)),
			)

			if err := p.Bootstrap(); err != nil {
				return errors.Wrap(err, "unable to bootstrap project")
			}

			// TODO: add prompt that the project bootstrapped

			return nil
		},
	}

	projName.RegisterString(cmd, &inputs.Name, "")
	projRoot.RegisterString(cmd, &inputs.Root, "")
	projCategory.RegisterString(cmd, &inputs.Category, "")
	projLayout.RegisterString(cmd, &inputs.Layout, "")

	return cmd
}

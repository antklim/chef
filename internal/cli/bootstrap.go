package cli

import (
	"fmt"
	"strings"

	"github.com/antklim/chef"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	projCategoryUnknown = "unknown"
	projCategoryCLI     = "cli"
	projCategoryPackage = "pkg"
	projCategoryService = "srv"
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
		Use:   "bootstrap",
		Args:  cobra.NoArgs,
		Short: "Bootstrap a new project",
		Long:  "Bootstrap a new project",
		Example: `chef bootstrap --name myproject
chef bootstrap --category [cli|pkg|srv] --name myproject
chef bootstrap -c [cli|pkg|srv] -n myproject --root /usr/local --layout chef.yml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projCategory := projCategoryFor(inputs.Category)

			if projCategory == projCategoryUnknown {
				return fmt.Errorf("unknown project category: %s", inputs.Category)
			}

			// TODO: init project structure and pass it to Init
			// project := project {
			// 	Name: inputs.Name,
			// 	...
			// }

			if err := chef.Init(inputs.Name); err != nil {
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

func projCategoryFor(v string) string {
	switch strings.ToLower(v) {
	case "cli":
		return projCategoryCLI
	case "pkg", "package":
		return projCategoryPackage
	case "srv", "service":
		return projCategoryService
	default:
		return projCategoryUnknown
	}
}

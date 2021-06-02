package project

import "strings"

type Category string

const (
	// CategoryUnknown represents unknown category of a project.
	CategoryUnknown Category = "unknown"
	// CategoryCLI represents CLI category of a project.
	CategoryCLI Category = "cli"
	// CategoryPackage represents package category of a project.
	CategoryPackage Category = "pkg"
	// CategoryService represents service category of a project.
	CategoryService Category = "srv"
)

func CategoryFor(v string) Category {
	switch strings.ToLower(v) {
	case "cli":
		return CategoryCLI
	case "pkg", "package":
		return CategoryPackage
	case "srv", "service":
		return CategoryService
	default:
		return CategoryUnknown
	}
}

func (c Category) IsUnknown() bool {
	return c == CategoryUnknown
}

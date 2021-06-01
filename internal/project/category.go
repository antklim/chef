package project

import "strings"

type Category string

const (
	CategoryUnknown Category = "unknown"
	CategoryCLI     Category = "cli"
	CategoryPackage Category = "pkg"
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

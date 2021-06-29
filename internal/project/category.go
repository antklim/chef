package project

import "strings"

const (
	categoryUnknown = "unknown"
	categoryCLI     = "cli"
	categoryPackage = "pkg"
	categoryService = "srv"
)

func category(v string) string {
	switch strings.ToLower(v) {
	case "cli":
		return categoryCLI
	case "pkg", "package":
		return categoryPackage
	case "srv", "service":
		return categoryService
	default:
		return categoryUnknown
	}
}

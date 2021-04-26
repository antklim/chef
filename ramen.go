package ramen

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// Project manager.
type Project struct {
	Taste  string
	Server string
}

var defaultProject = &Project{
	Taste:  "app",
	Server: "http",
}

// New project.
func New() *Project {
	return defaultProject
}

// Init initializes the project layout.
// TODO: make init a package level function with the default project.Init call.
func (p *Project) Init(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("project name required")
	}

	return os.Mkdir(name, 0755)
}

// Location returns project location for a given project name and root.
func Location(name, root string) (string, error) {
	wd, err := Root(root)
	if err != nil {
		return "", err
	}

	loc := path.Join(wd, name)

	fi, err := os.Stat(loc)
	if fi != nil {
		return "", fmt.Errorf("file or directory %s already exists", name)
	}

	return loc, nil
}

// Root validates provided project root.
func Root(name string) (string, error) {
	var err error

	if name = strings.TrimSpace(name); name == "" {
		name, err = os.Getwd()
	}

	if err != nil {
		return "", err
	}

	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}

	if !fi.IsDir() {
		return "", fmt.Errorf("%s is not a directory", name)
	}

	return name, nil
}

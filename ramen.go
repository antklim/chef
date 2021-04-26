package ramen

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// TODO: make taste and server enums
// TODO: make location and root private/internal
// TODO: make a list of layout components like cmd, internal, ...

// Project manager.
type Project struct {
	Name   string
	Root   string
	Taste  string
	Server string
}

var defaultProject = &Project{
	Name:   "ramen",
	Taste:  "app",
	Server: "http",
}

// New project.
// TODO: add options
// TODO: move name validation to option validation
func New() *Project {
	return defaultProject
}

// Init initializes the project layout.
// TODO: make init a package level function with the default project.Init call.
func (p *Project) Init(name, root string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("project name required")
	}

	loc, err := Location(name, root)
	if err != nil {
		return err
	}

	// TODO: refactor directory creation
	// 1. Making project root dir
	if err := os.Mkdir(loc, 0755); err != nil {
		return err
	}

	// 2. Making root/cmd
	if err := os.Mkdir(path.Join(loc, "cmd"), 0755); err != nil {
		return err
	}

	if _, err := os.Create(path.Join(loc, "cmd", "main.go")); err != nil {
		return err
	}

	// 3. Making root/internal
	if err := os.Mkdir(path.Join(loc, "internal"), 0755); err != nil {
		return err
	}

	if err := os.Mkdir(path.Join(loc, "internal", "app"), 0755); err != nil {
		return err
	}

	if err := os.Mkdir(path.Join(loc, "internal", "adapter"), 0755); err != nil {
		return err
	}

	if err := os.Mkdir(path.Join(loc, "internal", "provider"), 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(path.Join(loc, "internal", "server", "http"), 0755); err != nil {
		return err
	}

	if err := os.Mkdir(path.Join(loc, "test"), 0755); err != nil {
		return err
	}

	return nil
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

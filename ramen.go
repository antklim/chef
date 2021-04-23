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
}

// New project.
func New() *Project {
	return &Project{}
}

// Init initializes the project layout.
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

func Root(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return os.Getwd()
	}

	return name, nil
}

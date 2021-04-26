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
// TODO: read layout settings from yaml
// TODO: update location to location validation function

type LayoutDir int

const (
	DirCmd LayoutDir = iota
	DirInternal
	DirTest
	DirApp
	DirAdapter
	DirProvider
	DirServer
	DirHttp
)

var dirName = map[LayoutDir]string{
	DirCmd:      "cmd",
	DirInternal: "internal",
	DirTest:     "test",
	DirApp:      "app",
	DirAdapter:  "adapter",
	DirProvider: "provider",
	DirServer:   "server",
	DirHttp:     "http",
}

type node int

const (
	nodeDir node = iota
	nodeFile
)

type layoutNode struct {
	Name     string
	Type     node
	Children []layoutNode
}

var defaultLayout = []layoutNode{
	{
		Name: dirName[DirCmd],
		Children: []layoutNode{
			{Name: "main.go", Type: nodeFile},
		},
	},
	{
		Name: dirName[DirInternal],
		Children: []layoutNode{
			{Name: dirName[DirApp]},
			{Name: dirName[DirAdapter]},
			{Name: dirName[DirProvider]},
			{
				Name: dirName[DirServer],
				Children: []layoutNode{
					{Name: dirName[DirHttp]},
				},
			},
		},
	},
	{Name: dirName[DirTest]},
}

func layoutBuilder(root string, node layoutNode) error {
	o := path.Join(root, node.Name) // file system object, either file or directory

	switch node.Type {
	case nodeFile:
		f, err := os.Create(o)
		if err != nil {
			return err
		}
		return f.Chmod(0644)
	default:
		if err := os.Mkdir(o, 0755); err != nil {
			return err
		}

		for _, c := range node.Children {
			if err := layoutBuilder(o, c); err != nil {
				return err
			}
		}
	}

	return nil
}

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

	_, err := Location(name, root)
	if err != nil {
		return err
	}

	rl := layoutNode{
		Name:     name,
		Children: defaultLayout,
	}

	if err := layoutBuilder(root, rl); err != nil {
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

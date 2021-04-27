package ramen

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// TODO: read layout settings from yaml
// TODO: make location and root functions private/internal
// TODO: update location to location validation function

type ProjectTaste string

const (
	TasteApp ProjectTaste = "app"
	TastePkg ProjectTaste = "pkg"
)

type ProjectServer string

const (
	ServerHttp ProjectServer = "http"
)

type layoutDir int

const (
	dirCmd layoutDir = iota
	dirInternal
	dirTest
	dirApp
	dirAdapter
	dirProvider
	dirServer
	dirHttp
)

var dirName = map[layoutDir]string{
	dirCmd:      "cmd",
	dirInternal: "internal",
	dirTest:     "test",
	dirApp:      "app",
	dirAdapter:  "adapter",
	dirProvider: "provider",
	dirServer:   "server",
	dirHttp:     "http",
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
		Name: dirName[dirCmd],
		Children: []layoutNode{
			{Name: "main.go", Type: nodeFile},
		},
	},
	{
		Name: dirName[dirInternal],
		Children: []layoutNode{
			{Name: dirName[dirApp]},
			{Name: dirName[dirAdapter]},
			{Name: dirName[dirProvider]},
			{
				Name: dirName[dirServer],
				Children: []layoutNode{
					{Name: dirName[dirHttp]},
				},
			},
		},
	},
	{Name: dirName[dirTest]},
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
	Taste  ProjectTaste
	Server ProjectServer
}

func defaultProject(name string) Project {
	return Project{
		Name:   name,
		Taste:  TasteApp,
		Server: ServerHttp,
	}
}

// New project.
// TODO: add options
func New(name string) Project {
	name = strings.TrimSpace(name)
	return defaultProject(name)
}

// TODO: implement option validation
func (p *Project) Validate() error {
	if p.Name == "" {
		return errors.New("project name required: empty name provided")
	}

	return nil
}

// Init initializes the project layout.
// TODO: make init a package level function with the default project.Init call.
func (p *Project) Init(root string) error {
	_, err := Location(p.Name, root)
	if err != nil {
		return err
	}

	rl := layoutNode{
		Name:     p.Name,
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

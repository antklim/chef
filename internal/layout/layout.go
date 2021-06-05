package layout

import (
	"os"
	"path"
	"text/template"
)

// TODO: read layout settings from yaml
// TODO: test/build generated go code

// TODO: use http handler root template when bootstraping project
// TODO: use http handler template to add health endpoint (on bootstrap)
// TODO: make adding health endpoint on bootstrap optional
// TODO: update server main template

// TODO: support functionality of bring your own templates

// TODO: init project with go.mod

type layoutDir int

const (
	dirAdapter layoutDir = iota + 1
	dirApp
	dirCmd
	dirHandler
	dirHTTP
	dirInternal
	dirPkg
	dirProvider
	dirServer
	dirTest
)

var dirName = map[layoutDir]string{
	dirAdapter:  "adapter",
	dirApp:      "app",
	dirCmd:      "cmd",
	dirHandler:  "handler",
	dirHTTP:     "http",
	dirInternal: "internal",
	dirPkg:      "pkg",
	dirProvider: "provider",
	dirServer:   "server",
	dirTest:     "test",
}

type node int

const (
	nodeDir node = iota
	nodeFile
)

const (
	fperm = 0644
	dperm = 0755
)

// TODO: consider using the following interfaces
type nnode struct {
	Name        string
	Permissions int
}

type DirNode interface {
	Children() []nnode
}

type FileNode interface {
	Template() *template.Template
}

type Node struct {
	Name     string
	Type     node
	Children []Node
}

// Default project layout.
// TODO: make it private
// var defaultServiceLayout = []Node{
var Default = []Node{
	{Name: dirName[dirAdapter]},
	{Name: dirName[dirApp]},
	{
		Name: dirName[dirHandler],
		Children: []Node{
			{
				Name: dirName[dirHTTP],
				Children: []Node{
					{Name: "router.go", Type: nodeFile},
				},
			},
		},
	},
	{Name: dirName[dirProvider]},
	{
		Name: dirName[dirServer],
		Children: []Node{
			{
				Name: dirName[dirHTTP],
				Children: []Node{
					{Name: "server.go", Type: nodeFile},
				},
			},
		},
	},
	{Name: dirName[dirTest]},
	{Name: "main.go", Type: nodeFile},
}

func Builder(root string, n Node) error {
	o := path.Join(root, n.Name) // file system object, either file or directory

	switch n.Type {
	case nodeFile:
		f, err := os.Create(o)
		if err != nil {
			return err
		}
		defer f.Close()

		// TODO: refactor
		if n.Name == "main.go" {
			if err := srvMainTemplate.Execute(f, nil); err != nil {
				return err
			}
		}

		return f.Chmod(fperm)
	case nodeDir:
		fallthrough
	default:
		if err := os.Mkdir(o, dperm); err != nil {
			return err
		}

		for _, c := range n.Children {
			if err := Builder(o, c); err != nil {
				return err
			}
		}

		if len(n.Children) == 0 {
			if err := os.WriteFile(path.Join(o, ".gitkeep"), nil, fperm); err != nil {
				return err
			}
		}
	}

	return nil
}

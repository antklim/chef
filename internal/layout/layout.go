package layout

import (
	"os"
	"path"
)

type layoutDir int

const (
	dirCmd layoutDir = iota + 1
	dirInternal
	dirTest
	dirApp
	dirAdapter
	dirProvider
	dirServer
	dirHTTP
)

var dirName = map[layoutDir]string{
	dirCmd:      "cmd",
	dirInternal: "internal",
	dirTest:     "test",
	dirApp:      "app",
	dirAdapter:  "adapter",
	dirProvider: "provider",
	dirServer:   "server",
	dirHTTP:     "http",
}

type node int

const (
	nodeDir node = iota
	nodeFile
)

type Node struct {
	Name     string
	Type     node
	Children []Node
}

// Default project layout.
// TODO: make it private
var Default = []Node{
	{
		Name: dirName[dirCmd],
		Children: []Node{
			{Name: "main.go", Type: nodeFile},
		},
	},
	{
		Name: dirName[dirInternal],
		Children: []Node{
			{Name: dirName[dirApp]},
			{Name: dirName[dirAdapter]},
			{Name: dirName[dirProvider]},
			{
				Name: dirName[dirServer],
				Children: []Node{
					{Name: dirName[dirHTTP]},
				},
			},
		},
	},
	{Name: dirName[dirTest]},
}

func Builder(root string, n Node) error {
	o := path.Join(root, n.Name) // file system object, either file or directory

	switch n.Type {
	case nodeFile:
		f, err := os.Create(o)
		if err != nil {
			return err
		}
		return f.Chmod(0644) // nolint
	case nodeDir:
		fallthrough
	default:
		if err := os.Mkdir(o, 0755); err != nil {
			return err
		}

		for _, c := range n.Children {
			if err := Builder(o, c); err != nil {
				return err
			}
		}
	}

	return nil
}

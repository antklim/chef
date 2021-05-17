package layout

import (
	_ "embed" // required to be able to use go:embed
	"os"
	"path"
)

// TODO: read layout settings from yaml
// TODO: add main.go via go:embed
// TODO: use assets directory to store files used in project generation

type layoutDir int

const (
	dirCmd layoutDir = iota + 1
	dirAdapter
	dirApp
	dirHandler
	dirInternal
	dirPkg
	dirProvider
	dirTest
)

var dirName = map[layoutDir]string{
	dirCmd:      "cmd",
	dirAdapter:  "adapter",
	dirApp:      "app",
	dirHandler:  "handler",
	dirInternal: "internal",
	dirPkg:      "pkg",
	dirProvider: "provider",
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

//go:embed assets/.gitkeep
var gitkeep []byte

//go:embed assets/app/cmd/main.go
var appmain []byte

type Node struct {
	Name     string
	Type     node
	Children []Node
}

// Default project layout.
// TODO: make it private
var Default = []Node{
	{Name: dirName[dirAdapter]},
	{Name: dirName[dirApp]},
	{Name: dirName[dirHandler]},
	{Name: dirName[dirProvider]},
	{Name: dirName[dirTest]},
	{Name: "main.go", Type: nodeFile},
}

func Builder(root string, n Node) error {
	o := path.Join(root, n.Name) // file system object, either file or directory

	switch n.Type {
	case nodeFile:
		// TODO: refactor
		if n.Name == "main.go" {
			if err := os.WriteFile(o, appmain, fperm); err != nil {
				return err
			}
		} else {
			f, err := os.Create(o)
			if err != nil {
				return err
			}
			return f.Chmod(fperm)
		}
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
			if err := os.WriteFile(path.Join(o, ".gitkeep"), gitkeep, fperm); err != nil {
				return err
			}
		}
	}

	return nil
}

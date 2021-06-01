package layout

import (
	"os"
	"path"
)

// TODO: read layout settings from yaml
// TODO: test/build generated go code

// TODO: rename bootstrap command to boot

// TODO: use http handler root template when bootstraping project
// TODO: use http handler template to add health endpoint (on bootstrap)
// TODO: make adding health endpoint on bootstrap optional
// TODO: update server main template

// TODO: support functionality of bring your own templates

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

var gitkeep []byte

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
			if err := os.WriteFile(path.Join(o, ".gitkeep"), gitkeep, fperm); err != nil {
				return err
			}
		}
	}

	return nil
}

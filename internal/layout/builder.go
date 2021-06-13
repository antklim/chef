package layout

import (
	"io/fs"
	"os"
	"path"
	"text/template"
)

// TODO: init/register all possible layouts at start time.
// 			 When initing a project get one of the layouts and use it to scaffold
//       project structure.

// TODO: in imports replace chef/... with the project name

// TODO: read layout settings from yaml
// TODO: test/build generated go code

// TODO: use http handler template to add health endpoint (on bootstrap)
// TODO: make adding health endpoint on bootstrap optional

// TODO: support functionality of bring your own templates

// TODO: init project with go.mod

const (
	gitkeep = ".gitkeep"
)

const (
	fperm = 0644
	dperm = 0755
)

func Builder(root string, n Node) error {
	if nn, ok := n.(dirNode); ok {
		return buildDirNode(root, n, nn.SubNodes())
	}

	if nn, ok := n.(fileNode); ok {
		return buildFileNode(root, n, nn.Template())
	}

	return nil
}

func buildDirNode(root string, n Node, children []Node) error {
	o := path.Join(root, n.Name())

	if err := os.Mkdir(o, fs.FileMode(n.Permissions())); err != nil {
		return err
	}

	for _, c := range children {
		if err := Builder(o, c); err != nil {
			return err
		}
	}

	if len(children) == 0 {
		if err := os.WriteFile(path.Join(o, gitkeep), nil, fperm); err != nil {
			return err
		}
	}

	return nil
}

func buildFileNode(root string, n Node, t *template.Template) error {
	o := path.Join(root, n.Name())

	f, err := os.Create(o)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(f, nil); err != nil {
		return err
	}

	return f.Chmod(fs.FileMode(n.Permissions()))
}

// TODO: add options to define what subnodes layout to use
// TODO: deprecate root node

func RootNode(name string) Node {
	return newdnode(name, withSubNodes(defaultHTTPServiceLayout...))
}

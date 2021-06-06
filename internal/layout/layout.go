package layout

import (
	"io/fs"
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

const (
	dirAdapter  = "adapter"
	dirApp      = "app"
	dirCmd      = "cmd" // nolint
	dirHandler  = "handler"
	dirHTTP     = "http"
	dirInternal = "internal" // nolint
	dirPkg      = "pkg"      // nolint
	dirProvider = "provider"
	dirServer   = "server"
	dirTest     = "test"

	gitkeep = ".gitkeep"
)

const (
	fperm = 0644
	dperm = 0755
)

type dirNode interface {
	Children() []Node
}

type fileNode interface {
	Template() *template.Template
}

type Node interface {
	Name() string
	Permissions() uint32
}

func Builder(root string, n Node) error {
	if nn, ok := n.(dirNode); ok {
		return buildDirNode(root, n, nn.Children())
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

type dnode struct {
	name        string
	permissions uint32
	children    []Node
}

func (n dnode) Name() string {
	return n.name
}

func (n dnode) Permissions() uint32 {
	return n.permissions
}

func (n dnode) Children() []Node {
	return n.children
}

type fnode struct {
	name        string
	permissions uint32
	template    *template.Template
}

func (n fnode) Name() string {
	return n.name
}

func (n fnode) Permissions() uint32 {
	return n.permissions
}

func (n fnode) Template() *template.Template {
	return n.template
}

func RootNode(name string) Node {
	return dnode{
		name:        name,
		permissions: dperm,
		children:    defaultServiceLayout,
	}
}

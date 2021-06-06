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

type node struct {
	name        string
	permissions uint32
	children    []Node
}

func (n node) Name() string {
	return n.name
}

func (n node) Permissions() uint32 {
	return n.permissions
}

func (n node) Children() []Node {
	return n.children
}

func RootNode(name string) Node {
	return node{
		name:        name,
		permissions: dperm,
		children:    defaultServiceLayout,
	}
}

var defaultServiceLayout = []Node{
	node{
		name:        dirAdapter,
		permissions: dperm,
	},
	node{
		name:        dirApp,
		permissions: dperm,
	},
	node{
		name:        dirHandler,
		permissions: dperm,
		children: []Node{
			node{
				name:        dirHTTP,
				permissions: dperm,
				// TODO: add template
				// children: []Nnode{
				// 	httpRouter,
				// },
			},
		},
	},
	node{
		name:        dirProvider,
		permissions: dperm,
	},
	node{
		name:        dirServer,
		permissions: dperm,
		children: []Node{
			node{
				name:        dirHTTP,
				permissions: dperm,
				// TODO: add template
				// children: []Nnode{
				// 	httpServer,
				// },
			},
		},
	},
	node{
		name:        dirTest,
		permissions: dperm,
	},
	SrvMain,
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

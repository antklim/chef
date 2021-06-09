package layout

import (
	"io/fs"
	"os"
	"path"
	"text/template"
)

// TODO: in imports replace chef/... with the project name

// TODO: read layout settings from yaml
// TODO: test/build generated go code

// TODO: use http handler template to add health endpoint (on bootstrap)
// TODO: make adding health endpoint on bootstrap optional
// TODO: update server main template

// TODO: support functionality of bring your own templates

// TODO: init project with go.mod

const (
	gitkeep = ".gitkeep"
)

const (
	fperm = 0644
	dperm = 0755
)

type dirNode interface {
	SubNodes() []Node
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

type node struct {
	name        string
	permissions uint32
}

type dnode struct {
	node
	subnodes []Node
}

func newdnode(name string, opts ...dnodeoption) dnode {
	n := dnode{
		node: node{
			name:        name,
			permissions: dperm,
		},
	}

	for _, o := range opts {
		o.apply(&n)
	}

	return n
}

func (n dnode) Name() string {
	return n.name
}

func (n dnode) Permissions() uint32 {
	return n.permissions
}

func (n dnode) SubNodes() []Node {
	return n.subnodes
}

func (n *dnode) addSubNodes(sn []Node) {
	n.subnodes = append(n.subnodes, sn...)
}

type dnodeoption interface {
	apply(*dnode)
}

type dnodefopt struct {
	f func(*dnode)
}

func (f *dnodefopt) apply(n *dnode) {
	f.f(n)
}

func newdnodefopt(f func(*dnode)) *dnodefopt {
	return &dnodefopt{f}
}

func withSubNodes(sn ...Node) dnodeoption {
	return newdnodefopt(func(n *dnode) {
		n.subnodes = sn
	})
}

func withPermissions(p uint32) dnodeoption {
	return newdnodefopt(func(n *dnode) {
		n.permissions = p
	})
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

// TODO: add options to define what subnodes layout to use

func RootNode(name string) Node {
	return dnode{
		node: node{
			name:        name,
			permissions: dperm,
		},
		subnodes: defaultHTTPServiceLayout,
	}
}

package layout

import (
	"io/fs"
	"os"
	"path"
	"text/template"
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
	node
	template *template.Template
}

func (n fnode) Name() string {
	return n.name
}

func (n fnode) Permissions() uint32 {
	return n.permissions
}

// Build executes node template and writes it to a file to a provided directory.
func (n fnode) Build(dir string) error {
	if n.template == nil {
		return nil
	}

	o := path.Join(dir, n.Name())

	f, err := os.Create(o)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := n.template.Execute(f, nil); err != nil {
		return err
	}

	return f.Chmod(fs.FileMode(n.Permissions()))
}

func (n fnode) Template() *template.Template {
	return n.template
}

// TODO: add newfnode
// TODO: add fnode write method

package layout

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"text/template"
)

// TODO: add flag wasbuild to nodes to support rebuild operation
// TODO: Dnode and Fnode should be interfaces

const (
	fperm fs.FileMode = 0644
	dperm fs.FileMode = 0755
)

var (
	errNilTemplate = errors.New("node template is nil")
)

type Node interface {
	Name() string
	Permissions() fs.FileMode
	Build(loc, mod string) error
}

type node struct {
	name        string
	permissions fs.FileMode
}

// Dnode describes directory nodes.
// Use NewDnode to initialize Dnode.
type Dnode struct {
	node
	subnodes []Node
}

func NewDnode(name string, opts ...DnodeOption) Dnode {
	n := Dnode{
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

func (n Dnode) Name() string {
	return n.name
}

func (n Dnode) Permissions() fs.FileMode {
	return n.permissions
}

func (n Dnode) Build(loc, mod string) error {
	o := path.Join(loc, n.Name())

	if err := os.Mkdir(o, n.Permissions()); err != nil {
		return err
	}

	for _, sn := range n.subnodes {
		if err := sn.Build(o, mod); err != nil {
			return err
		}
	}

	return nil
}

func (n Dnode) SubNodes() []Node {
	return n.subnodes
}

func (n Dnode) GetSubNode(name string) Node {
	return findByName(n.subnodes, name)
}

func (n *Dnode) AddSubNode(newNode Node) error {
	if subnode := findByName(n.subnodes, newNode.Name()); subnode != nil {
		return fmt.Errorf("subnode %s already exists", newNode.Name())
	}
	n.subnodes = append(n.subnodes, newNode)
	return nil
}

func findByName(nodes []Node, n string) Node {
	for _, node := range nodes {
		if node.Name() == n {
			return node
		}
	}
	return nil
}

type DnodeOption interface {
	apply(*Dnode)
}

type dnodefopt struct {
	f func(*Dnode)
}

func (f *dnodefopt) apply(n *Dnode) {
	f.f(n)
}

func newdnodefopt(f func(*Dnode)) *dnodefopt {
	return &dnodefopt{f}
}

func WithSubNodes(sn ...Node) DnodeOption {
	return newdnodefopt(func(n *Dnode) {
		n.subnodes = sn
	})
}

func WithDperm(p fs.FileMode) DnodeOption {
	return newdnodefopt(func(n *Dnode) {
		n.permissions = p
	})
}

// Fnode describes file nodes.
// Use NewFnode to initialize Fnode.
type Fnode struct {
	node
	template *template.Template
}

func NewFnode(name string, opts ...FnodeOption) Fnode {
	n := Fnode{
		node: node{
			name:        name,
			permissions: fperm,
		},
	}

	for _, o := range opts {
		o.apply(&n)
	}

	return n
}

func (n Fnode) Name() string {
	return n.name
}

func (n Fnode) Permissions() fs.FileMode {
	return n.permissions
}

// Build executes node template and writes it to a file to a provided location.
func (n Fnode) Build(loc, mod string) error {
	if n.template == nil {
		return errNilTemplate
	}

	o := path.Join(loc, n.Name())

	f, err := os.Create(o)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := n.wbuild(f, mod); err != nil {
		return err
	}

	return f.Chmod(n.Permissions())
}

func (n Fnode) wbuild(w io.Writer, mod string) error {
	data := struct {
		Module string
	}{
		Module: mod,
	}

	return n.template.Execute(w, data)
}

func (n Fnode) Template() *template.Template {
	return n.template
}

type FnodeOption interface {
	apply(*Fnode)
}

type fnodefopt struct {
	f func(*Fnode)
}

func (f *fnodefopt) apply(n *Fnode) {
	f.f(n)
}

func newfnodefopt(f func(*Fnode)) *fnodefopt {
	return &fnodefopt{f}
}

func WithFperm(p fs.FileMode) FnodeOption {
	return newfnodefopt(func(n *Fnode) {
		n.permissions = p
	})
}

// WithNewTemplate adds node template with template name tn and template string
// ts.
func WithNewTemplate(tn, ts string) FnodeOption {
	return newfnodefopt(func(n *Fnode) {
		n.template = template.Must(template.New(tn).Parse(ts))
	})
}

// WithTemplate adds node template t.
func WithTemplate(t *template.Template) FnodeOption {
	return newfnodefopt(func(n *Fnode) {
		n.template = t
	})
}

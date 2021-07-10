package layout

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"text/template"
)

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

func (n *Dnode) AddSubNodes(sn []Node) {
	n.subnodes = append(n.subnodes, sn...)
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

type fnode struct {
	node
	template *template.Template
}

func newfnode(name string, opts ...fnodeoption) fnode {
	n := fnode{
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

func (n fnode) Name() string {
	return n.name
}

func (n fnode) Permissions() fs.FileMode {
	return n.permissions
}

// Build executes node template and writes it to a file to a provided location.
func (n fnode) Build(loc, mod string) error {
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

func (n fnode) wbuild(w io.Writer, mod string) error {
	data := struct {
		Module string
	}{
		Module: mod,
	}

	return n.template.Execute(w, data)
}

func (n fnode) Template() *template.Template {
	return n.template
}

type fnodeoption interface {
	apply(*fnode)
}

type fnodefopt struct {
	f func(*fnode)
}

func (f *fnodefopt) apply(n *fnode) {
	f.f(n)
}

func newfnodefopt(f func(*fnode)) *fnodefopt {
	return &fnodefopt{f}
}

func withFperm(p fs.FileMode) fnodeoption {
	return newfnodefopt(func(n *fnode) {
		n.permissions = p
	})
}

// withNewTemplate adds node template with template name tn and template string
// ts.
func withNewTemplate(tn, ts string) fnodeoption {
	return newfnodefopt(func(n *fnode) {
		n.template = template.Must(template.New(tn).Parse(ts))
	})
}

// withNewTemplate adds node template t.
func withTemplate(t *template.Template) fnodeoption {
	return newfnodefopt(func(n *fnode) {
		n.template = t
	})
}

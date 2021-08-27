package node

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/pkg/errors"
)

// TODO (feat): consider to store location in node and remove it from Build

const (
	fperm fs.FileMode = 0644
	dperm fs.FileMode = 0755
)

var (
	errNilTemplate = errors.New("node template is nil")
)

// Adder is the interface that wraps node Add method.
//
// Add adds Node to a collection of subnodes.
// It returns an error if Node could not be added to the collection of subnodes.
type Adder interface {
	Add(Node) error
}

// Getter is the interface that wraps node Get method.
//
// Get searches node in the provided location in the collection of subnodes.
// It returns nil if no nodes found.
type Getter interface {
	Get(string) Node
}

// Node interface defines layout node functionality.
type Node interface {
	// Name returns a node name.
	Name() string
	// Build executes node build.
	Build(loc string, data interface{}) error
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

// NewDnode creates a new directory node.
func NewDnode(name string, opts ...DnodeOption) *Dnode {
	n := &Dnode{
		node: node{
			name:        name,
			permissions: dperm,
		},
	}

	for _, o := range opts {
		o.apply(n)
	}

	return n
}

// Name returns a node name.
func (n *Dnode) Name() string {
	return n.name
}

// Build creates a directory in file system recursively builds all subnodes.
//
// When subnode build fails the process stops and the error is returned.
// Node directory is not deleted in case of build failure.
func (n *Dnode) Build(loc string, data interface{}) error {
	o := path.Join(loc, n.Name())

	if err := os.Mkdir(o, n.permissions); err != nil {
		return err
	}

	for _, sn := range n.subnodes {
		if err := sn.Build(o, data); err != nil {
			return errors.Wrapf(err, "failed to build subnode %q", sn.Name())
		}
	}

	return nil
}

// Nodes returns a list of subnodes.
func (n *Dnode) Nodes() []Node {
	return n.subnodes
}

// Get returns the first node found by name in all subnodes of the node, or nil
// if no node found.
func (n *Dnode) Get(name string) Node {
	return findByName(n.subnodes, name)
}

// Add adds a new node to a list of subnodes.
//
// When subnode list already has a node with the same name the error returned.
func (n *Dnode) Add(newNode Node) error {
	if subnode := findByName(n.subnodes, newNode.Name()); subnode != nil {
		return fmt.Errorf("node %q already exists", newNode.Name())
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

// DnodeOption sets directory node options.
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

// WithSubNodes returns an DnodeOption that sets directory node subnodes.
func WithSubNodes(sn ...Node) DnodeOption {
	return newdnodefopt(func(n *Dnode) {
		n.subnodes = sn
	})
}

// WithDperm returns an DnodeOption that sets directory permissions.
func WithDperm(p fs.FileMode) DnodeOption {
	return newdnodefopt(func(n *Dnode) {
		n.permissions = p
	})
}

// Fnode describes file nodes.
type Fnode struct {
	node
	template *template.Template
}

// NewFnode creates a new file node.
func NewFnode(name string, opts ...FnodeOption) *Fnode {
	n := &Fnode{
		node: node{
			name:        name,
			permissions: fperm,
		},
	}

	for _, o := range opts {
		o.apply(n)
	}

	return n
}

// Name returns a node name.
func (n *Fnode) Name() string {
	return n.name
}

// Build executes node template and writes it to a file to a provided location.
func (n *Fnode) Build(loc string, data interface{}) error {
	if n.template == nil {
		return errNilTemplate
	}

	o := path.Join(loc, n.Name())

	f, err := os.Create(o)
	if err != nil {
		return err
	}
	if err := n.wbuild(f, data); err != nil {
		f.Close()
		// Remove created file. It's a clean up, thus ignore errors here.
		os.Remove(o)
		return errors.Wrap(err, "failed to execute template")
	}
	defer f.Close()

	return f.Chmod(n.permissions)
}

func (n *Fnode) wbuild(w io.Writer, data interface{}) error {
	return n.template.Execute(w, data)
}

// FnodeOption sets file node options.
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

// WithFperm returns an FnodeOption that sets file permissions.
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

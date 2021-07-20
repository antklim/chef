package layout

import (
	"fmt"
	"path"
	"strings"
)

const Root = "."

type Dir interface {
	Add(n Node) error
	Get(string) Node
	Nodes() []Node
}

type Layout struct {
	root   Dir
	schema string
}

// New creates a new layout with schema s and nodes n.
func New(s string, nodes ...Node) Layout {
	root := NewDnode(Root, WithSubNodes(nodes...))
	return Layout{
		root:   root,
		schema: s,
	}
}

func (l Layout) Nodes() []Node {
	return l.root.Nodes()
}

// Add adds a node to a layout location.
func (l *Layout) Add(n Node, loc string) error {
	if node := l.Get(n.Name(), loc); node != nil {
		return fmt.Errorf("node %s already exists at '%s'", n.Name(), loc)
	}

	if loc == Root {
		return l.root.Add(n)
	}

	locNode := l.Get(path.Base(loc), path.Dir(loc))
	if locNode == nil {
		return fmt.Errorf("path '%s' not found in layout", loc)
	}

	locDir, ok := locNode.(Dir)
	if !ok {
		return fmt.Errorf("node '%s' does not support adding subnodes", loc)
	}

	return locDir.Add(n)
}

func (l Layout) Schema() string {
	return l.schema
}

func (l Layout) Build(loc, mod string) error {
	for _, n := range l.root.Nodes() {
		if err := n.Build(loc, mod); err != nil {
			return err
		}
	}
	return nil
}

// Get returns a node with the given name at a location.
func (l Layout) Get(node, loc string) Node {
	if loc == Root {
		return l.root.Get(node)
	}

	dirs := strings.Split(loc, "/")
	d := l.root

	for _, dir := range dirs {
		n := d.Get(dir)
		if n == nil {
			return nil
		}

		dnode, ok := n.(Dir)
		if !ok {
			return nil
		}
		d = dnode
	}

	return d.Get(node)
}

const (
	// ServiceLayout an abstract service layout name.
	ServiceLayout = "srv"
	// HTTPServiceLayout an http service layout name.
	HTTPServiceLayout = "srv_http"
)

var (
	// m is a map from schema to layout.
	m = make(map[string]Layout)
)

// Register registers the layout to the layouts map. l.Schema will be
// used as the schema registered with this layout.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe.
func Register(l Layout) {
	m[l.Schema()] = l
}

// Get returns the layout registered with the given schema.
//
// If no layout is registered with the schema, Nil layout will be returned.
func Get(name string) *Layout {
	if l, ok := m[name]; ok {
		return &l
	}
	return nil
}

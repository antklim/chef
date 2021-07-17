package layout

import (
	"fmt"
	"strings"
)

const Root = ""

type Layout struct {
	root   Dnode
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
	return l.root.SubNodes()
}

// Add adds a node to a layout location.
func (l *Layout) Add(n Node, loc string) error {
	if l.Has(n.Name(), loc) {
		return fmt.Errorf("node %s already exists at '%s'", n.Name(), loc)
	}

	if loc == Root {
		return l.root.AddSubNode(n)
	}

	return nil
}

func (l Layout) Schema() string {
	return l.schema
}

func (l Layout) Build(loc, mod string) error {
	for _, n := range l.root.SubNodes() {
		if err := n.Build(loc, mod); err != nil {
			return err
		}
	}
	return nil
}

// Has returns true if layout has a node at a location.
func (l Layout) Has(node, loc string) bool {
	if loc == Root {
		_, ok := find(l.root.SubNodes(), node)
		return ok
	}

	dirs := strings.Split(loc, "/")
	d := l.root

	for _, dir := range dirs {
		n := d.GetSubNode(dir)
		if n == nil {
			return false
		}

		dnode, ok := n.(Dnode)
		if !ok {
			return false
		}
		d = dnode
	}

	n := d.GetSubNode(node)

	return n != nil
}

func find(nodes []Node, node string) (Node, bool) {
	for _, n := range nodes {
		if n.Name() == node {
			return n, true
		}
	}
	return nil, false
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

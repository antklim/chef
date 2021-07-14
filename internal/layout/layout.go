package layout

import (
	"fmt"
	"strings"
)

const Root = ""

type Layout struct {
	nodes  []Node
	schema string
}

// New creates a new layout with schema s and nodes n.
func New(s string, nodes ...Node) Layout {
	return Layout{
		schema: s,
		nodes:  nodes,
	}
}

func (l Layout) Nodes() []Node {
	return l.nodes
}

// Add adds a node to a layout location.
func (l *Layout) Add(n Node, loc string) error {
	if l.Has(n.Name(), loc) {
		return fmt.Errorf("node %s already exists at '%s'", n.Name(), loc)
	}

	if loc == Root {
		l.nodes = append(l.nodes, n)
		return nil
	}

	return nil
}

func (l Layout) Schema() string {
	return l.schema
}

func (l Layout) Build(loc, mod string) error {
	for _, n := range l.nodes {
		if err := n.Build(loc, mod); err != nil {
			return err
		}
	}
	return nil
}

// Has returns true if layout has a node at a location.
func (l Layout) Has(node, loc string) bool {
	if loc == "" {
		_, ok := find(l.nodes, node)
		return ok
	}

	dirs := strings.Split(loc, "/")
	d := NewDnode("_", WithSubNodes(l.nodes...))

	for _, dir := range dirs {
		n, ok := find(d.SubNodes(), dir)
		if !ok {
			return false
		}

		if d, ok = n.(Dnode); !ok {
			return false
		}
	}

	_, ok := find(d.SubNodes(), node)

	return ok
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

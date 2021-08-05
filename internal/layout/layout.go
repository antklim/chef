package layout

import (
	"fmt"
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
	rootNode := NewDnode(Root, WithSubNodes(nodes...))
	root := NewDnode("", WithSubNodes(rootNode))
	return Layout{
		root:   root,
		schema: s,
	}
}

// AddNode adds a node to a layout location.
func (l *Layout) AddNode(n Node, loc string) error {
	locNode := l.FindNode(loc)
	if locNode == nil {
		return fmt.Errorf("node %q not found in layout", loc)
	}

	locDir, ok := locNode.(Dir)
	if !ok {
		return fmt.Errorf("node %q does not support adding subnodes", loc)
	}

	if node := locDir.Get(n.Name()); node != nil {
		return fmt.Errorf("node %q already has subnode %q", loc, n.Name())
	}

	return locDir.Add(n)
}

func (l Layout) Schema() string {
	return l.schema
}

func (l Layout) Build(loc, mod string) error {
	root := l.rootDir()
	for _, n := range root.Nodes() {
		if err := n.Build(loc, mod); err != nil {
			return err
		}
	}
	return nil
}

// TODO: format comment bellow for better documentation help

// FindNode returns a node at provided location.
// For example:
// - find("server/http") returns directory node associated with "server/http" location
// - find("server/http/handler.go") returns file node associated with the handler.go
// - find(".") returns root node
// - find("") returns nil when no associated node found
func (l Layout) FindNode(loc string) Node {
	locs := splitPath(loc)
	node := l.root

	for _, l := range locs[:len(locs)-1] {
		n := node.Get(l)
		if n == nil {
			return nil
		}

		dnode, ok := n.(Dir)
		if !ok {
			return nil
		}
		node = dnode
	}

	return node.Get(locs[len(locs)-1])
}

func (l Layout) rootDir() Dir {
	rootNode := l.root.Get(Root)
	if rootNode == nil {
		return nil
	}

	dir, ok := rootNode.(Dir)
	if !ok {
		return nil
	}
	return dir
}

func splitPath(loc string) []string {
	a := strings.Split(loc, "/")
	if a[0] != "." {
		a = append([]string{"."}, a...)
	}
	return a
}

// TODO: deprecate layout registry or move it to the lower level package
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

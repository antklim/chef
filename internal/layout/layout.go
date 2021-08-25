package layout

import (
	"fmt"
	"strings"

	"github.com/antklim/chef/internal/layout/node"
)

// A Root is a root node of the layout.
const Root = "."

// Dir interface defines directory node functionality.
type Dir interface {
	Add(n node.Node) error
	Get(string) node.Node
	Nodes() []node.Node
}

// A Layout defines project layout.
type Layout struct {
	root Dir
}

// New creates a new layout with nodes.
func New(nodes ...node.Node) *Layout {
	rootNode := node.NewDnode(Root, node.WithSubNodes(nodes...))
	root := node.NewDnode("", node.WithSubNodes(rootNode))
	return &Layout{root: root}
}

// AddNode adds a node to the location in the layout.
func (l *Layout) AddNode(n node.Node, loc string) error {
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

// Build recursively builds all nodes in layout.
func (l *Layout) Build(loc, mod string) error {
	data := struct {
		Module string
	}{
		Module: mod,
	}

	root := l.rootDir()
	for _, n := range root.Nodes() {
		if err := n.Build(loc, data); err != nil {
			return err
		}
	}
	return nil
}

// TODO (ref): format comment bellow for better documentation help

// FindNode returns a node at the provided location in the layout.
// For example:
// - find("server/http") returns directory node associated with "server/http" location
// - find("server/http/handler.go") returns file node associated with the handler.go
// - find(".") returns root node
// - find("") returns nil when no associated node found
func (l *Layout) FindNode(loc string) node.Node {
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

func (l *Layout) rootDir() Dir {
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
	if a[0] != Root {
		a = append([]string{Root}, a...)
	}
	return a
}

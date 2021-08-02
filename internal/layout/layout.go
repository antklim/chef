package layout

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

var (
	errComponentTemplateNil = errors.New("component template is nil")
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

// TODO: refactor root node in a way that Get and findNode does not consider root as a separate use-case.

// New creates a new layout with schema s and nodes n.
func New(s string, nodes ...Node) Layout {
	root := NewDnode(Root, WithSubNodes(nodes...))
	return Layout{
		root:   root,
		schema: s,
	}
}

// TODO: Consider making it private

// AddNode adds a node to a layout location.
func (l *Layout) AddNode(n Node, loc string) error {
	if node := l.GetNode(n.Name(), loc); node != nil {
		return fmt.Errorf("node %s already exists at %q", n.Name(), loc)
	}

	if loc == Root {
		return l.root.Add(n)
	}

	locNode := l.GetNode(path.Base(loc), path.Dir(loc))
	if locNode == nil {
		return fmt.Errorf("path %q not found in layout", loc)
	}

	locDir, ok := locNode.(Dir)
	if !ok {
		return fmt.Errorf("node %q does not support adding subnodes", loc)
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

// GetNode returns a node with the given name at a location.
func (l Layout) GetNode(node, loc string) Node {
	if loc == Root {
		return l.root.Get(node)
	}

	n := l.findNode(loc)
	d, ok := n.(Dir)
	if !ok {
		return nil
	}
	return d.Get(node)
}

// Find returns a node referenced by location.
func (l Layout) Find(loc string) Node {
	return nil
}

// AddComponent adds component node to the layout.
// func (l *Layout) AddComponent(componentName, nodeName string) error {
// 	component, ok := l.components[componentName]
// 	if !ok {
// 		return fmt.Errorf("unknown component %q", componentName)
// 	}

// 	if node := l.GetNode(nodeName, component.loc); node != nil {
// 		return fmt.Errorf("%s %q already exists", componentName, nodeName)
// 	}

// 	node := NewFnode(nodeName, WithTemplate(component.template))
// 	return l.AddNode(node, component.loc)
// }

// TODO: refactor - unify Get and findNode
// TODO: format comment bellow for better documentation help

// find searches for a node at provided location.
// For example:
// - find("server/http") returns directory node associated with "server/http" location
// - find("server/http/handler.go") returns file node associated with the handler.go
// - find(".") returns nil for root location
// - find("") returns nil when no associated node found
func (l Layout) findNode(loc string) Node {
	dirs := strings.Split(loc, "/")
	d := l.root

	for _, dir := range dirs[:len(dirs)-1] {
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

	return d.Get(dirs[len(dirs)-1])
}

// TODO: consider deprecation of layout registry
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

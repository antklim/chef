package layout_test

import (
	"errors"
	"strings"

	"github.com/antklim/chef/internal/layout/node"
)

type testNode struct {
	name        string
	buildCalled bool
	buildError  error
	loc         string
}

func newTestNode(name string) *testNode {
	return &testNode{name: name}
}

func (n *testNode) Build(loc string, data interface{}) error {
	n.buildCalled = true
	n.loc = loc
	if strings.HasPrefix(loc, "/error") {
		n.buildError = errors.New("node build error")
		return n.buildError
	}
	return nil
}

func (n *testNode) Name() string {
	return n.name
}

func (n *testNode) WasBuild() bool {
	return n.buildCalled == true
}

func (n *testNode) BuiltAt() string {
	return n.loc
}

var _ node.Node = (*testNode)(nil)

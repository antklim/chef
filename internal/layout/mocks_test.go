package layout_test

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/antklim/chef/internal/layout"
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

func (*testNode) Permissions() fs.FileMode {
	return 0400
}

func (n *testNode) WasBuild() bool {
	return n.buildCalled == true
}

func (n *testNode) BuiltAt() string {
	return n.loc
}

var _ layout.Node = (*testNode)(nil)

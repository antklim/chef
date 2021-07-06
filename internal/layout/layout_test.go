package layout_test

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testNode struct {
	buildCalled bool
	buildError  error
	loc         string
}

func (n *testNode) Build(loc, mod string) error {
	n.buildCalled = true
	n.loc = loc
	if strings.HasPrefix(loc, "/error") {
		n.buildError = errors.New("node build error")
		return n.buildError
	}
	return nil
}

func (testNode) Name() string {
	return "testNode"
}

func (testNode) Permissions() fs.FileMode {
	return 0400
}

func (n testNode) WasBuild() bool {
	return n.buildCalled == true
}

func (n testNode) BuiltAt() string {
	return n.loc
}

var _ layout.Node = (*testNode)(nil)

func TestNewLayout(t *testing.T) {
	testCases := []struct {
		desc   string
		schema string
		nodes  []layout.Node
	}{
		{
			desc:   "creates layout with defined schema",
			schema: "testLayout",
		},
		{
			desc:   "creates layout with nodes",
			schema: "testLayoutWNodes",
			nodes:  []layout.Node{&testNode{}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l := layout.New(tC.schema, tC.nodes)
			assert.Equal(t, tC.schema, l.Schema())
			assert.Equal(t, tC.nodes, l.Nodes())
		})
	}
}

func TestLayoutBuild(t *testing.T) {
	testCases := []struct {
		desc   string
		name   string
		node   *testNode
		loc    string
		assert func(*testing.T, error)
	}{
		{
			desc: "builds layout nodes",
			name: "successBuild",
			node: &testNode{},
			loc:  "/tmp/foo/bar",
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			desc: "fails to build layout when node build fails",
			name: "errorBuild",
			node: &testNode{},
			loc:  "/error/bar",
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "node build error")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l := layout.New(tC.name, []layout.Node{tC.node})
			assert.False(t, tC.node.WasBuild())
			err := l.Build(tC.loc, "module_name")
			tC.assert(t, err)
			assert.True(t, tC.node.WasBuild())
			assert.Equal(t, tC.loc, tC.node.BuiltAt())
		})
	}
}

func TestLayoutRegistry(t *testing.T) {
	t.Run("get returns nil when layout not registered", func(t *testing.T) {
		l := layout.Get("foo")
		assert.Nil(t, l)
	})

	t.Run("get returns layout by schema", func(t *testing.T) {
		tl := layout.New("testLayout", nil)
		layout.Register(tl)
		l := layout.Get("testLayout")
		assert.Equal(t, tl, *l)
	})
}

func TestLayoutInit(t *testing.T) {
	t.Run("registers predefined layouts", func(t *testing.T) {
		defs := []string{layout.ServiceLayout, layout.HTTPServiceLayout}
		for _, s := range defs {
			l := layout.Get(s)
			assert.NotNil(t, l)
		}
	})
}

func TestLayoutHas(t *testing.T) {
	t.Skip("not implemented")

	// Prepare l Layout
	testCases := []struct {
		node     string
		loc      string
		expected bool
	}{
		{
			// is true for top level node
		},
		{
			// is true for nested node
		},
		{
			// is true for file node
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("is %t for node %s at %s", tC.expected, tC.node, tC.loc), func(t *testing.T) {

		})
	}
}

func TestDefaultLayouts(t *testing.T) {
	t.Run("service layout has correct nodes", func(t *testing.T) {
		t.Skip("not implemented")
		// l := layout.Get(layout.ServiceLayout)
	})

	t.Run("http service layout has correct nodes", func(t *testing.T) {
		t.Skip("not implemented")
		// l := layout.Get(layout.HTTPServiceLayout)
	})
}

// func TestLayoutAddEndpoint(t *testing.T) {
// 	l := layout.Get(layout.HTTPServiceLayout)
// 	l.AddEndpoint(&testNode{})
// 	assert.True(l.Has("testNode", "handler/http"))
// }

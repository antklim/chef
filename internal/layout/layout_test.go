package layout_test

import (
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testNode struct {
	buildCalled bool
	loc         string
}

func (n *testNode) Build(loc string) error {
	n.buildCalled = true
	n.loc = loc
	return nil
}

func (testNode) Name() string {
	return "testNode"
}

func (testNode) Permissions() uint32 {
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
	n := &testNode{}
	l := layout.New("testBuild", []layout.Node{n})
	assert.False(t, n.WasBuild())

	loc := "/tmp/foo/bar"
	err := l.Build(loc)
	require.NoError(t, err)
	assert.True(t, n.WasBuild())
	assert.Equal(t, loc, n.BuiltAt())
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
		assert.Equal(t, tl, l)
	})

	t.Run("has predefined layouts", func(t *testing.T) {
		defs := []string{"srv", "srv_http"}
		for _, s := range defs {
			l := layout.Get(s)
			assert.NotNil(t, l)
		}
	})
}

package layout_test

import (
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
)

func TestDnodeGetSubNode(t *testing.T) {
	fnode := layout.NewFnode("file.txt")
	dnode := layout.NewDnode("dnode", layout.WithSubNodes(fnode))

	testCases := []struct {
		desc     string
		name     string
		expected layout.Node
	}{
		{
			desc:     "returns sub node by name",
			name:     "file.txt",
			expected: fnode,
		},
		{
			desc:     "returns nil when node not found",
			name:     "other-file.txt",
			expected: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			n := dnode.GetSubNode(tC.name)
			assert.Equal(t, tC.expected, n)
		})
	}
}

func TestDnodeAddSubNode(t *testing.T) {
	fnode := layout.NewFnode("file.txt")
	dnode := layout.NewDnode("dnode", layout.WithSubNodes(fnode))

	t.Run("returns an error when existing sub node has same name as the new", func(t *testing.T) {
		subnodesBefore := len(dnode.SubNodes())

		newNode := layout.NewDnode("file.txt")
		err := dnode.AddSubNode(newNode)
		assert.EqualError(t, err, "node file.txt already exists")

		subnodesAfter := len(dnode.SubNodes())
		assert.Equal(t, subnodesBefore, subnodesAfter)
	})

	t.Run("adds a new subnode", func(t *testing.T) {
		subnodesBefore := len(dnode.SubNodes())

		newNode := layout.NewFnode("file2.txt")
		err := dnode.AddSubNode(newNode)
		assert.NoError(t, err)

		subnodesAfter := len(dnode.SubNodes())
		assert.Equal(t, subnodesBefore+1, subnodesAfter)
	})
}

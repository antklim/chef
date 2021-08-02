package layout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNode(t *testing.T) {
	handlerNode := NewFnode("handler.go")
	httpNode := NewDnode("http", WithSubNodes(handlerNode))
	serverNode := NewDnode("server", WithSubNodes(httpNode))
	l := New("layout", serverNode)

	testCases := []struct {
		desc     string
		loc      string
		expected Node
	}{
		{
			desc:     "returns server/http node",
			loc:      "server/http",
			expected: httpNode,
		},
		{
			desc:     "returns handler node",
			loc:      "server/http/handler.go",
			expected: handlerNode,
		},
		{
			desc: "returns nil when node does not exist",
			loc:  "server/http/other.go",
		},
		{
			desc: "returns nil for root",
			loc:  Root,
		},
		{
			desc: "returns nil when location does not exist",
			loc:  "",
		},
		{
			desc: "returns nil when location does not exist",
			loc:  "server/grpc",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			node := l.findNode(tC.loc)
			assert.Equal(t, tC.expected, node)
		})
	}
}

package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLayoutFactory(t *testing.T) {
	testCases := []struct {
		desc     string
		category string
		serever  string
	}{
		{
			desc:     "returns nil for unknown category",
			category: "foo",
		},
		{
			desc:     "returns nil for service category and unknown server",
			category: "srv",
			serever:  "bar",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			lf := layoutFactory(category(tC.category), server(tC.serever))
			assert.Nil(t, lf)
		})
	}
}

func TestServiceLayoutFactory(t *testing.T) {
	lf := layoutFactory(category("srv"), server(""))
	assert.NotNil(t, lf)
	l := lf.makeLayout()
	assert.NotNil(t, l)

	expectedNodes := []string{"adapter", "app", "handler", "provider", "server", "test"}
	for _, n := range expectedNodes {
		node := l.FindNode(n)
		assert.NotNil(t, node)
	}
}

func TestHTTPServiceLayoutFactory(t *testing.T) {
	lf := layoutFactory(category("service"), server("http"))
	assert.NotNil(t, lf)
	l := lf.makeLayout()
	assert.NotNil(t, l)

	expectedNodes := []string{"adapter", "app", "handler", "provider", "server", "test", "main.go"}
	for _, n := range expectedNodes {
		node := l.FindNode(n)
		assert.NotNil(t, node)
	}
}

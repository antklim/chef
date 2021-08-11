package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentsFactory(t *testing.T) {
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
			desc:     "returns nil for service category without server",
			category: "srv",
		},
		{
			desc:     "returns nil for service category and unknown server",
			category: "srv",
			serever:  "bar",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			lf := componentsFactory(category(tC.category), server(tC.serever))
			assert.Nil(t, lf)
		})
	}
}

func TestHTTPServiceComponentsFactory(t *testing.T) {
	f := componentsFactory(category("service"), server("http"))
	assert.NotNil(t, f)
	c := f.makeComponents()
	assert.NotNil(t, c)

	expectedComponents := []string{"http_handler"}
	for _, v := range expectedComponents {
		assert.Contains(t, c, v)
	}
	assert.Len(t, c, len(expectedComponents))
}

package layout_test

import (
	"testing"

	"github.com/antklim/chef/internal/layout"
	"github.com/stretchr/testify/assert"
)

func TestLayoutSelector(t *testing.T) {
	t.Run("returns unknown category error when unknown category provided", func(t *testing.T) {
		nodes, err := layout.Selector("", "")
		assert.EqualError(t, err, "unknown layout category")
		assert.Nil(t, nodes)
	})

	t.Run("returns default service layout for service category", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("returns http service layout for service category and http server", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("returns unknown server error for service category and unknown server", func(t *testing.T) {
		nodes, err := layout.Selector(layout.CategoryService, "")
		assert.EqualError(t, err, "unknown server")
		assert.Nil(t, nodes)
	})
}

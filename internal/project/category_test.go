package project

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCategory(t *testing.T) {
	testCases := []struct {
		v        string
		expected string
	}{
		{
			v:        "cLi",
			expected: categoryCLI,
		},
		{
			v:        "pkg",
			expected: categoryPackage,
		},
		{
			v:        "package",
			expected: categoryPackage,
		},
		{
			v:        "SRV",
			expected: categoryService,
		},
		{
			v:        "service",
			expected: categoryService,
		},
		{
			v:        "foo",
			expected: categoryUnknown,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("returns %s category when %s provided", tC.expected, tC.v), func(t *testing.T) {
			actual := category(tC.v)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

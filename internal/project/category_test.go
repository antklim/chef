package project_test

import (
	"fmt"
	"testing"

	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
)

func TestNewCategory(t *testing.T) {
	testCases := []struct {
		v        string
		expected project.Category
	}{
		{
			v:        "cli",
			expected: project.CategoryCLI,
		},
		{
			v:        "pkg",
			expected: project.CategoryPackage,
		},
		{
			v:        "package",
			expected: project.CategoryPackage,
		},
		{
			v:        "srv",
			expected: project.CategoryService,
		},
		{
			v:        "service",
			expected: project.CategoryService,
		},
		{
			v:        "foo",
			expected: project.CategoryUnknown,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("returns %s category when %s provided", tC.expected, tC.v), func(t *testing.T) {
			actual := project.NewCategory(tC.v)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestCategoryIsUknown(t *testing.T) {
	testCases := []struct {
		v        project.Category
		expected bool
	}{
		{
			v:        project.CategoryCLI,
			expected: false,
		},
		{
			v:        project.CategoryPackage,
			expected: false,
		},
		{
			v:        project.CategoryService,
			expected: false,
		},
		{
			v:        project.CategoryUnknown,
			expected: true,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("returns %t for %s category", tC.expected, tC.v), func(t *testing.T) {
			assert.Equal(t, tC.expected, tC.v.IsUnknown())
		})
	}
}

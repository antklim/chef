package project_test

import (
	"fmt"
	"testing"

	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
)

func TestCategoryFor(t *testing.T) {
	testCases := []struct {
		desc     string
		v        string
		expected project.Category
	}{
		{
			desc:     "returns CLI category when %s provided",
			v:        "cli",
			expected: project.CategoryCLI,
		},
		{
			desc:     "returns Package category when %s provided",
			v:        "pkg",
			expected: project.CategoryPackage,
		},
		{
			desc:     "returns Package category when %s provided",
			v:        "package",
			expected: project.CategoryPackage,
		},
		{
			desc:     "returns Service category when %s provided",
			v:        "srv",
			expected: project.CategoryService,
		},
		{
			desc:     "returns Service category when %s provided",
			v:        "service",
			expected: project.CategoryService,
		},
		{
			desc:     "returns Unknown category when %s provided",
			v:        "foo",
			expected: project.CategoryUnknown,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf(tC.desc, tC.v), func(t *testing.T) {
			actual := project.CategoryFor(tC.v)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestIsUknown(t *testing.T) {
	testCases := []struct {
		desc     string
		v        project.Category
		expected bool
	}{
		{
			desc:     "returns false for CLI category",
			v:        project.CategoryCLI,
			expected: false,
		},
		{
			desc:     "returns false for Package category",
			v:        project.CategoryPackage,
			expected: false,
		},
		{
			desc:     "returns false for Service category",
			v:        project.CategoryService,
			expected: false,
		},
		{
			desc:     "returns true for Unknown category",
			v:        project.CategoryUnknown,
			expected: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.expected, tC.v.IsUnknown())
		})
	}
}

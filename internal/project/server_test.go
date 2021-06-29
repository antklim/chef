package project_test

import (
	"fmt"
	"testing"

	"github.com/antklim/chef/internal/project"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	testCases := []struct {
		v        string
		expected project.Server
	}{
		{
			v:        "",
			expected: project.ServerNone,
		},
		{
			v:        "Http",
			expected: project.ServerHTTP,
		},
		{
			v:        "grpc",
			expected: project.ServerGRPC,
		},
		{
			v:        "foo",
			expected: project.ServerUnknown,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("returns %s server when %s provided", tC.expected, tC.v), func(t *testing.T) {
			actual := project.NewServer(tC.v)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestServerIsUknown(t *testing.T) {
	testCases := []struct {
		v        project.Server
		expected bool
	}{
		{
			v:        project.ServerGRPC,
			expected: false,
		},
		{
			v:        project.ServerHTTP,
			expected: false,
		},
		{
			v:        project.ServerNone,
			expected: false,
		},
		{
			v:        project.ServerUnknown,
			expected: true,
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("returns %t for %s server", tC.expected, tC.v), func(t *testing.T) {
			assert.Equal(t, tC.expected, tC.v.IsUnknown())
		})
	}
}

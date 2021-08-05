package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	// returns error when trying to make project of unknown category
	// for category srv: makes service project srvProject, srvProject implements Project interface
	// testCases := []struct {
	// 	category string
	// 	expected interface{}
	// }{
	// 	{
	// 		category: "srv",
	// 		expected: srvProject(nil),
	// 	},
	// 	{
	// 		category: "service",
	// 		expected: srvProject(nil),
	// 	},
	// }
	// for _, tC := range testCases {
	// 	t.Run(fmt.Sprintf("makes project for %q category", tC.expected, tC.category), func(t *testing.T) {
	// 		p, err := MakeProject(tC.category)
	// 		require.NoError(t, err)
	// 		assert.IsType(t, tC.expected, p)
	// 		assert.Implements(t, IProject, p)
	// 	})
	// }

	t.Run("returns error when asked to make a project of unknown category", func(t *testing.T) {
		_, err := MakeProject("foo")
		assert.EqualError(t, err, `unknown project category "foo"`)
	})
}

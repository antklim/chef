package project

import "github.com/antklim/chef/internal/layout"

type Blueprint interface {
	Layout() *layout.Layout
	Components() []Component
}

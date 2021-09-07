package api

import (
	"text/template"

	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/layout/node"
	"github.com/antklim/chef/internal/project"
)

var testTmpl = template.Must(template.New("test").Parse("package foo"))

// ProjectFactory generates a project with default properties. By default it
// returns inited project with one registered component "http_handler".
func ProjectFactory(opts ...project.Option) (*project.Project, error) {
	dopts := defaultProjectOptions()
	oopts := make([]project.Option, 0, len(opts)+len(dopts))
	oopts = append(oopts, dopts...)
	oopts = append(oopts, opts...)

	p := project.New("cheftestapi", oopts...)
	if err := p.Init(); err != nil {
		return nil, err
	}

	c := project.NewComponent("http_handler", "handler", "HTTP Handler", testTmpl)
	if err := p.RegisterComponent(c); err != nil {
		return nil, err
	}

	return p, nil
}

func defaultProjectOptions() []project.Option {
	l := layout.New(node.NewDnode("handler"))
	opts := []project.Option{project.WithLayout(l)}
	return opts
}

// FindComponent finds and returns the first component at which f(c) is true.
// If component not found it returns nil.
func FindComponent(components []project.Component, f func(project.Component) bool) *project.Component {
	for _, c := range components {
		if f(c) {
			return &c
		}
	}
	return nil
}

package project

import (
	"path"
	"text/template"

	templ "github.com/antklim/chef/internal/project/template"
)

const (
	httpHandler = "http_handler"
)

type Component struct {
	name string
	loc  string
	desc string
	tmpl *template.Template
}

func NewComponent(name, loc, desc string, tmpl *template.Template) Component {
	return Component{
		name: name,
		loc:  loc,
		desc: desc,
		tmpl: tmpl,
	}
}

type componentsMaker interface {
	makeComponents() map[string]Component
}

func componentsFactory(category, server string) componentsMaker {
	if category == categoryService && server == serverHTTP {
		return httpServiceComponets{}
	}
	return nil
}

type httpServiceComponets struct{}

func (httpServiceComponets) makeComponents() map[string]Component {
	c := make(map[string]Component)
	c[httpHandler] = Component{
		loc:  path.Join(dirHandler, dirHTTP),
		name: httpHandler,
		desc: "HTTP handler",
		tmpl: templ.Get(templ.HTTPEndpoint),
	}
	return c
}

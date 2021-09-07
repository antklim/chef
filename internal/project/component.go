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
	Name string
	Loc  string
	Desc string
	Tmpl *template.Template
}

func NewComponent(name, loc, desc string, tmpl *template.Template) Component {
	return Component{
		Name: name,
		Loc:  loc,
		Desc: desc,
		Tmpl: tmpl,
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
		Name: httpHandler,
		Loc:  path.Join(dirHandler, dirHTTP),
		Desc: "HTTP handler",
		Tmpl: templ.Get(templ.HTTPEndpoint),
	}
	return c
}

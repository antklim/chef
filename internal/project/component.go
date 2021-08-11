package project

import (
	"path"
	"text/template"

	templ "github.com/antklim/chef/internal/project/template"
)

const (
	httpHandler = "http_handler"
)

type component struct {
	loc      string
	name     string
	template *template.Template
}

type componentsMaker interface {
	makeComponents() map[string]component
}

func componentsFactory(category, server string) componentsMaker {
	if category == categoryService && server == serverHTTP {
		return httpServiceComponets{}
	}
	return nil
}

type httpServiceComponets struct{}

func (httpServiceComponets) makeComponents() map[string]component {
	c := make(map[string]component)
	c[httpHandler] = component{
		loc:      path.Join(dirHandler, dirHTTP),
		name:     httpHandler,
		template: templ.Get(templ.HTTPEndpoint),
	}
	return c
}

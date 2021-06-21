package template

import "text/template"

const (
	HTTPRouter = "http_router"
)

var rootTemplate = template.New("__chef_root__")

// Get returns the template registered with the given name.
func Get(name string) *template.Template {
	return rootTemplate.Lookup(name)
}

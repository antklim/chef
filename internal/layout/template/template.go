package template

import "text/template"

const (
	// HTTPHandler an http handler template name.
	HTTPHandler = "http_handler"
	// HTTPRouter an http router template name.
	HTTPRouter = "http_router"
	// HTTPServer an http server template name.
	HTTPServer = "http_server"
)

var rootTemplate = template.New("__chef_root__")

// Get returns the template registered with the given name.
func Get(name string) *template.Template {
	return rootTemplate.Lookup(name)
}

package template

import "text/template"

const (
	// HTTPEndpoint an http endpoint template name.
	HTTPEndpoint = "http_endpoint"
	// HTTPRouter an http router template name.
	HTTPRouter = "http_router"
	// HTTPServer an http server template name.
	HTTPServer = "http_server"
	// HTTPService an http service template name.
	HTTPService = "http_service"
)

var rootTemplate = template.New("__chef_root__")

// Get returns the template registered with the given name.
func Get(name string) *template.Template {
	return rootTemplate.Lookup(name)
}

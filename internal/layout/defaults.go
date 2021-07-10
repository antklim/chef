package layout

import (
	"fmt"

	"github.com/antklim/chef/internal/layout/template"
)

const (
	dirAdapter  = "adapter"
	dirApp      = "app"
	dirHandler  = "handler"
	dirHTTP     = "http"
	dirServer   = "server"
	dirProvider = "provider"
	dirTest     = "test"
)

func init() { // nolint:gochecknoinits
	Register(New(ServiceLayout, serviceNodes()))
	Register(New(HTTPServiceLayout, httpServiceNodes()))
}

func httpEndpoint(name string) fnode {
	return newfnode(fmt.Sprintf("%s.go", name), withTemplate(template.Get(template.HTTPEndpoint)))
}

func serviceNodes() []Node {
	return []Node{
		NewDnode(dirAdapter),
		NewDnode(dirApp),
		NewDnode(dirHandler),
		NewDnode(dirProvider),
		NewDnode(dirServer),
		NewDnode(dirTest),
	}
}

func httpServiceNodes() []Node {
	httpRouter := newfnode("router.go", withTemplate(template.Get(template.HTTPRouter)))
	httpHandlerNode := NewDnode(dirHTTP, WithSubNodes(httpRouter))
	httpServer := newfnode("server.go", withTemplate(template.Get(template.HTTPServer)))
	httpServerNode := NewDnode(dirHTTP, WithSubNodes(httpServer))
	httpSrvMain := newfnode("main.go", withTemplate(template.Get(template.HTTPService)))

	return []Node{
		NewDnode(dirAdapter),
		NewDnode(dirApp),
		NewDnode(dirHandler, WithSubNodes(httpHandlerNode)),
		NewDnode(dirProvider),
		NewDnode(dirServer, WithSubNodes(httpServerNode)),
		NewDnode(dirTest),
		httpSrvMain,
	}
}

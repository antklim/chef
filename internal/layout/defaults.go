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
		NewDNode(dirAdapter),
		NewDNode(dirApp),
		NewDNode(dirHandler),
		NewDNode(dirProvider),
		NewDNode(dirServer),
		NewDNode(dirTest),
	}
}

func httpServiceNodes() []Node {
	httpRouter := newfnode("router.go", withTemplate(template.Get(template.HTTPRouter)))
	httpHandlerNode := NewDNode(dirHTTP, withSubNodes(httpRouter))
	httpServer := newfnode("server.go", withTemplate(template.Get(template.HTTPServer)))
	httpServerNode := NewDNode(dirHTTP, withSubNodes(httpServer))
	httpSrvMain := newfnode("main.go", withTemplate(template.Get(template.HTTPService)))

	return []Node{
		NewDNode(dirAdapter),
		NewDNode(dirApp),
		NewDNode(dirHandler, withSubNodes(httpHandlerNode)),
		NewDNode(dirProvider),
		NewDNode(dirServer, withSubNodes(httpServerNode)),
		NewDNode(dirTest),
		httpSrvMain,
	}
}

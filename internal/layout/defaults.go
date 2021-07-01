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

func httpHandler(name string) fnode {
	return newfnode(fmt.Sprintf("%s.go", name), withTemplate(template.Get(template.HTTPHandler)))
}

func serviceNodes() []Node {
	return []Node{
		newdnode(dirAdapter),
		newdnode(dirApp),
		newdnode(dirHandler),
		newdnode(dirProvider),
		newdnode(dirServer),
		newdnode(dirTest),
	}
}

func httpServiceNodes() []Node {
	httpRouter := newfnode("router.go", withTemplate(template.Get(template.HTTPRouter)))
	httpHandlerNode := newdnode(dirHTTP, withSubNodes(httpRouter))
	httpServer := newfnode("server.go", withTemplate(template.Get(template.HTTPServer)))
	httpServerNode := newdnode(dirHTTP, withSubNodes(httpServer))
	httpSrvMain := newfnode("main.go", withTemplate(template.Get(template.HTTPService)))

	return []Node{
		newdnode(dirAdapter),
		newdnode(dirApp),
		newdnode(dirHandler, withSubNodes(httpHandlerNode)),
		newdnode(dirProvider),
		newdnode(dirServer, withSubNodes(httpServerNode)),
		newdnode(dirTest),
		httpSrvMain,
	}
}

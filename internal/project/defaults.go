package project

import (
	"fmt"

	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/project/template"
)

// TODO: validate executed templates content

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
	layout.Register(layout.New(ServiceLayout, serviceNodes()...))
	layout.Register(layout.New(HTTPServiceLayout, httpServiceNodes()...))
}

func httpEndpoint(name string) *layout.Fnode {
	return layout.NewFnode(fmt.Sprintf("%s.go", name), layout.WithTemplate(template.Get(template.HTTPEndpoint)))
}

func serviceNodes() []layout.Node {
	return []layout.Node{
		layout.NewDnode(dirAdapter),
		layout.NewDnode(dirApp),
		layout.NewDnode(dirHandler),
		layout.NewDnode(dirProvider),
		layout.NewDnode(dirServer),
		layout.NewDnode(dirTest),
	}
}

func httpServiceNodes() []layout.Node {
	httpRouter := layout.NewFnode("router.go", layout.WithTemplate(template.Get(template.HTTPRouter)))
	httpHandlerNode := layout.NewDnode(dirHTTP, layout.WithSubNodes(httpRouter))
	httpServer := layout.NewFnode("server.go", layout.WithTemplate(template.Get(template.HTTPServer)))
	httpServerNode := layout.NewDnode(dirHTTP, layout.WithSubNodes(httpServer))
	httpSrvMain := layout.NewFnode("main.go", layout.WithTemplate(template.Get(template.HTTPService)))

	return []layout.Node{
		layout.NewDnode(dirAdapter),
		layout.NewDnode(dirApp),
		layout.NewDnode(dirHandler, layout.WithSubNodes(httpHandlerNode)),
		layout.NewDnode(dirProvider),
		layout.NewDnode(dirServer, layout.WithSubNodes(httpServerNode)),
		layout.NewDnode(dirTest),
		httpSrvMain,
	}
}

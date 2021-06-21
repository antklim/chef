package layout

import (
	"fmt"

	"github.com/antklim/chef/internal/layout/template"
)

const (
	dirAdapter = "adapter"
	dirApp     = "app"
	// dirCmd      = "cmd"
	dirHandler = "handler"
	dirHTTP    = "http"
	// dirInternal = "internal"
	// dirPkg      = "pkg"
	dirServer   = "server"
	dirProvider = "provider"
	dirTest     = "test"
)

func init() { // nolint:gochecknoinits
	Register(serviceLayout())
	Register(httpServiceLayout())
}

func httpHandler(name string) fnode {
	return newfnode(fmt.Sprintf("%s.go", name), withTemplate(template.Get(template.HTTPHandler)))
}

func serviceLayout() Layout {
	serviceNodes := []Node{
		newdnode(dirAdapter),
		newdnode(dirApp),
		newdnode(dirHandler),
		newdnode(dirProvider),
		newdnode(dirServer),
		newdnode(dirTest),
	}

	return New(ServiceLayout, serviceNodes)
}

func httpServiceLayout() Layout {
	httpRouter := newfnode("router.go", withTemplate(template.Get(template.HTTPRouter)))
	httpHandlerNode := newdnode(dirHTTP, withSubNodes(httpRouter))
	httpServer := newfnode("server.go", withTemplate(template.Get(template.HTTPServer)))
	httpServerNode := newdnode(dirHTTP, withSubNodes(httpServer))
	httpSrvMain := newfnode("main.go", withTemplate(template.Get(template.HTTPService)))

	httpServiceNodes := []Node{
		newdnode(dirAdapter),
		newdnode(dirApp),
		newdnode(dirHandler, withSubNodes(httpHandlerNode)),
		newdnode(dirProvider),
		newdnode(dirServer, withSubNodes(httpServerNode)),
		newdnode(dirTest),
		httpSrvMain,
	}

	return New(HTTPServiceLayout, httpServiceNodes)
}

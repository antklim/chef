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
	Register(serviceLayout)
	Register(httpServiceLayout)
}

func httpHandler(name string) fnode {
	return newfnode(fmt.Sprintf("%s.go", name), withTemplate(template.Get(template.HTTPHandler)))
}

var httpRouter = newfnode("router.go", withTemplate(template.Get(template.HTTPRouter)))
var httpServer = newfnode("server.go", withTemplate(template.Get(template.HTTPServer)))
var httpSrvMain = newfnode("main.go", withTemplate(template.Get(template.HTTPService)))

var serviceNodes = []Node{
	newdnode(dirAdapter),
	newdnode(dirApp),
	newdnode(dirHandler),
	newdnode(dirProvider),
	newdnode(dirServer),
	newdnode(dirTest),
}

// TODO: srv should be public constant
var serviceLayout = New("srv", serviceNodes)

var httpHandlerNode = newdnode(dirHTTP, withSubNodes(httpRouter))
var httpServerNode = newdnode(dirHTTP, withSubNodes(httpServer))

var httpServiceNodes = []Node{
	newdnode(dirAdapter),
	newdnode(dirApp),
	newdnode(dirHandler, withSubNodes(httpHandlerNode)),
	newdnode(dirProvider),
	newdnode(dirServer, withSubNodes(httpServerNode)),
	newdnode(dirTest),
	httpSrvMain,
}

// TODO: srv_http should be public constant
var httpServiceLayout = New("srv_http", httpServiceNodes)

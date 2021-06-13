package layout

func init() { // nolint:gochecknoinits
	Register(serviceLayout{})
	Register(httpServiceLayout{})
}

var serviceNodes = []Node{
	newdnode(dirAdapter),
	newdnode(dirApp),
	newdnode(dirHandler),
	newdnode(dirProvider),
	newdnode(dirServer),
	newdnode(dirTest),
}

type serviceLayout struct{}

func (serviceLayout) Nodes() []Node {
	return serviceNodes
}

func (serviceLayout) Schema() string {
	return "srv" // TODO: should be public constant
}

var httpHandlerNode = newdnode(dirHTTP, withSubNodes(httpRouter))
var httpServerNode = newdnode(dirHTTP, withSubNodes(httpServer))

var defaultHTTPServiceLayout = []Node{
	newdnode(dirAdapter),
	newdnode(dirApp),
	newdnode(dirHandler, withSubNodes(httpHandlerNode)),
	newdnode(dirProvider),
	newdnode(dirServer, withSubNodes(httpServerNode)),
	newdnode(dirTest),
	srvMain,
}

type httpServiceLayout struct{}

func (httpServiceLayout) Nodes() []Node {
	return defaultHTTPServiceLayout
}

func (httpServiceLayout) Schema() string {
	return "srv_http" // TODO: should be public constant
}

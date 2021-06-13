package layout

func init() { // nolint:gochecknoinits
	Register(serviceLayout{})
	Register(httpServiceLayout{})
}

var defaultServiceLayout = []Node{
	newdnode(dirAdapter),
	newdnode(dirApp),
	newdnode(dirHandler),
	newdnode(dirProvider),
	newdnode(dirServer),
	newdnode(dirTest),
}

type serviceLayout struct{}

func (serviceLayout) Nodes() []Node {
	return defaultServiceLayout
}

func (serviceLayout) Schema() string {
	return "srv"
}

var defaultHTTPHandlerLayout = newdnode(dirHTTP, withSubNodes(httpRouter))
var defaultHTTPServerLayout = newdnode(dirHTTP, withSubNodes(httpServer))

var defaultHTTPServiceLayout = []Node{
	newdnode(dirAdapter),
	newdnode(dirApp),
	newdnode(dirHandler, withSubNodes(defaultHTTPHandlerLayout)),
	newdnode(dirProvider),
	newdnode(dirServer, withSubNodes(defaultHTTPServerLayout)),
	newdnode(dirTest),
	srvMain,
}

type httpServiceLayout struct{}

func (httpServiceLayout) Nodes() []Node {
	return defaultHTTPServiceLayout
}

func (httpServiceLayout) Schema() string {
	return "srv_http"
}

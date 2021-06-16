package layout

func init() { // nolint:gochecknoinits
	Register(serviceLayout)
	Register(httpServiceLayout)
}

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
	srvMain,
}

// TODO: srv_http should be public constant
var httpServiceLayout = New("srv_http", httpServiceNodes)

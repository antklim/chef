package layout

var defaultServiceLayout = []Node{ // nolint
	newdnode(dirAdapter),
	newdnode(dirApp),
	newdnode(dirHandler),
	newdnode(dirProvider),
	newdnode(dirServer),
	newdnode(dirTest),
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

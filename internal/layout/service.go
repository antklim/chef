package layout

var defaultServiceLayout = []Node{
	dnode{
		name:        dirAdapter,
		permissions: dperm,
	},
	dnode{
		name:        dirApp,
		permissions: dperm,
	},
	dnode{
		name:        dirHandler,
		permissions: dperm,
		children: []Node{
			dnode{
				name:        dirHTTP,
				permissions: dperm,
				children: []Node{
					httpRouter,
				},
			},
		},
	},
	dnode{
		name:        dirProvider,
		permissions: dperm,
	},
	dnode{
		name:        dirServer,
		permissions: dperm,
		children: []Node{
			dnode{
				name:        dirHTTP,
				permissions: dperm,
				children: []Node{
					httpServer,
				},
			},
		},
	},
	dnode{
		name:        dirTest,
		permissions: dperm,
	},
	srvMain,
}

/*
var defaultHttpServiceLayout = []Node{
	// dirAdapter
	// dirApp
	// dirHandler(WithHttpRouter, WithHttpHandler(healthHandler))
	// dirProvider
	// dirServer(WithHttpServer)
	// dirTest
	// srvMain
}
*/

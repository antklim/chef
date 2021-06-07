package layout

var defaultServiceLayout = []Node{ // nolint
	dirAdapter,
	dirApp,
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
	dirProvider,
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
	dirTest,
	srvMain,
}

var defaultHTTPServiceLayout = []Node{
	dirAdapter,
	dirApp,
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
	dirProvider,
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
	dirTest,
	srvMain,
}

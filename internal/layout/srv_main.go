package layout

import "text/template"

const _srvMainTemplate = `package main

import (
	server "chef/server/http"
)

func main() {
	server.Start()
}
`

var srvMain = fnode{
	name:        "main.go",
	permissions: fperm,
	template:    template.Must(template.New("srv_main").Parse(_srvMainTemplate)),
}

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
				// TODO: add template
				// children: []Nnode{
				// 	httpRouter,
				// },
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
				// TODO: add template
				// children: []Nnode{
				// 	httpServer,
				// },
			},
		},
	},
	dnode{
		name:        dirTest,
		permissions: dperm,
	},
	srvMain,
}

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
	node{
		name:        dirAdapter,
		permissions: dperm,
	},
	node{
		name:        dirApp,
		permissions: dperm,
	},
	node{
		name:        dirHandler,
		permissions: dperm,
		children: []Node{
			node{
				name:        dirHTTP,
				permissions: dperm,
				// TODO: add template
				// children: []Nnode{
				// 	httpRouter,
				// },
			},
		},
	},
	node{
		name:        dirProvider,
		permissions: dperm,
	},
	node{
		name:        dirServer,
		permissions: dperm,
		children: []Node{
			node{
				name:        dirHTTP,
				permissions: dperm,
				// TODO: add template
				// children: []Nnode{
				// 	httpServer,
				// },
			},
		},
	},
	node{
		name:        dirTest,
		permissions: dperm,
	},
	srvMain,
}

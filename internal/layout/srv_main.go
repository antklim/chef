package layout

import "text/template"

// TODO: in imports replace chef/... with the project name

const srvMainTemplate = `package main

import (
	server "{{ .Module }}/server/http"
)

func main() {
	server.Start()
}
`

var srvMain = fnode{
	name:        "main.go",
	permissions: fperm,
	template:    template.Must(template.New("srv_main").Parse(srvMainTemplate)),
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

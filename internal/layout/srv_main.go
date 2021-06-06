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

type srvMainNode struct {
	name        string
	permissions uint32
	template    *template.Template
}

func (n srvMainNode) Name() string {
	return n.name
}

func (n srvMainNode) Permissions() uint32 {
	return n.permissions
}

func (n srvMainNode) Template() *template.Template {
	return n.template
}

var SrvMain = srvMainNode{
	name:        "main.go",
	permissions: fperm,
	template:    template.Must(template.New("srv_main").Parse(_srvMainTemplate)),
}

// type serverNode struct {
// 	name        string
// 	permissions uint32
// 	children    []Nnode
// }

// func (n serverNode) Name() string {
// 	return n.name
// }

// func (n serverNode) Permissions() uint32 {
// 	return n.permissions
// }

// func (n serverNode) Children() []Nnode {
// 	return n.children
// }

// var httpServer = serverNode{
// 	name:        "http",
// 	permissions: dperm,
// 	children:    []Nnode{srvMain},
// }

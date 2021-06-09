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

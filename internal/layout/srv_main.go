package layout

import "text/template"

// TODO: in imports replace chef/... with the project name

var srvMainTemplate = template.Must(template.New("srv_main").Parse(`package main

import (
	server "{{ .Module }}/server/http"
)

func main() {
	server.Start()
}
`))

var srvMain = newfnode("main.go", withTemplate(srvMainTemplate))

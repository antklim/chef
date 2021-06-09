package layout

import "text/template"

const httpRouterTemplate = `package http

import "net/http"

var router = http.NewServeMux()

func Mux() *http.ServeMux {
	return router
}
`

var httpRouter = fnode{
	node: node{
		name:        "router.go",
		permissions: fperm,
	},
	template: template.Must(template.New("http_router").Parse(httpRouterTemplate)),
}

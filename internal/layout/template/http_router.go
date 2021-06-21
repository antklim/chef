package template

import "text/template"

var _ = template.Must(rootTemplate.New("http_router").Parse(`package http

import "net/http"

var router = http.NewServeMux()

func Mux() *http.ServeMux {
	return router
}
`))

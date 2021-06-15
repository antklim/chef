package layout

import "text/template"

var httpRouterTemplate = template.Must(template.New("http_router").Parse(`package http

import "net/http"

var router = http.NewServeMux()

func Mux() *http.ServeMux {
	return router
}
`))

var httpRouter = newfnode("router.go", withTemplate(httpRouterTemplate))

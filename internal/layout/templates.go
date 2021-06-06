package layout

import "text/template"

// TODO: in imports replace chef/... with the project name

const _httpServerTemplate = `package http

import (
	handler "chef/handler/http"
	"log"
	"net/http"
)

const defaultAddress = ":8080"

func Start() {
	s := &http.Server{
		Addr:    defaultAddress,
		Handler: handler.Mux(),
	}

	log.Printf("service listening at %s", defaultAddress)
	log.Fatalf("service stopped: %v", s.ListenAndServe())
}
`

var httpServerTemplate = template.Must(template.New("http_server").Parse(_httpServerTemplate)) // nolint

const _httpRouterTemplate = `package http

import "net/http"

var router = http.NewServeMux()

func Mux() *http.ServeMux {
	return router
}
`

var httpRouterTemplate = template.Must(template.New("http_router").Parse(_httpRouterTemplate)) // nolint

const _httpRouteTemplate = `package http

import (
	"fmt"
	"net/http"
)

const {{ .Name }}Route = {{ .Path }}

func init() {
	router.Handle({{ .Name }}Route, {{ .Name }}Handler())
}

func {{ .Name }}Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
}
`

var httpRouteTemplate = template.Must(template.New("http_route").Parse(_httpRouteTemplate)) // nolint

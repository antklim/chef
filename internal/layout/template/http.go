package template

import "text/template"

var _ = template.Must(rootTemplate.New(HTTPRoute).Parse(`package http

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
`))

var _ = template.Must(rootTemplate.New(HTTPRouter).Parse(`package http

import "net/http"

var router = http.NewServeMux()

func Mux() *http.ServeMux {
	return router
}
`))

var _ = template.Must(rootTemplate.New(HTTPServer).Parse(`package http

import (
	handler "{{ .Module }}/handler/http"
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
`))

var _ = template.Must(rootTemplate.New(HTTPService).Parse(`package main

import (
	server "{{ .Module }}/server/http"
)

func main() {
	server.Start()
}
`))

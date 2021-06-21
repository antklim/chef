package template

import "text/template"

var _ = template.Must(rootTemplate.New(HTTPRouter).Parse(`package http

import "net/http"

var router = http.NewServeMux()

func Mux() *http.ServeMux {
	return router
}
`))

// TODO: in imports replace chef/... with the project name

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

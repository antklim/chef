package layout

import "text/template"

// TODO: in imports replace chef/... with the project name

var httpServerTemplate = template.Must(template.New("http_server").Parse(`package http

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

var httpServer = newfnode("server.go", withTemplate(httpServerTemplate))

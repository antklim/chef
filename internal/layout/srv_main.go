package layout

import "text/template"

const _srvMainTemplate = `package main

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	s *http.Server
}

func NewServer() *Server {
	mux := http.NewServeMux()

	r, h := handler()
	mux.Handle(r, h)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &Server{s: s}
}

func (s *Server) Start() error {
	return s.s.ListenAndServe()
}

func handler() (string, http.Handler) {
	route := "/health"
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	}
	return route, http.HandlerFunc(h)
}

func main() {
	srv := NewServer()
	log.Println("Starting a service...")
	log.Fatal(srv.Start())
}
`

var srvMainTemplate = template.Must(template.New("srv_main").Parse(_srvMainTemplate))

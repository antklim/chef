package layout

// const _httpServerTemplate = `
// package http

// import (
// 	handler {{ .HTTPHandlerPkg }}
// 	"log"
// 	"net/http"
// )

// const defaultAddress = ":8080"

// func Start() {
// 	mux := http.NewServeMux()

// 	handler.Muxer(mux)

// 	s := &http.Server{
// 		Addr:    defaultAddress,
// 		Handler: mux,
// 	}

// 	log.Printf("service listening at %s", defaultAddress)
// 	log.Fatalf("service stopped: %v", s.ListenAndServe())
// }
// `

// var httpServerTemplate = template.Must(template.New("http_server").Parse(_httpServerTemplate))

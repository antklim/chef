package layout

import (
	"fmt"
	"text/template"
)

const httpHandlerTemplate = `package http

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

func httpHandler(name string) fnode {
	return fnode{
		node: node{
			name:        fmt.Sprintf("%s.go", name),
			permissions: fperm,
		},
		template: template.Must(template.New("http_handler").Parse(httpHandlerTemplate)),
	}
}

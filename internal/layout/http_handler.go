package layout

import (
	"fmt"
	"text/template"
)

var httpHandlerTemplate = template.Must(template.New("http_handler").Parse(`package http

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

func httpHandler(name string) fnode {
	return newfnode(fmt.Sprintf("%s.go", name), withTemplate(httpHandlerTemplate))
}

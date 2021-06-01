package layout

import "text/template"

const _httpHandlerTemplate = `
package http

import (
	"fmt"
	"net/http"
)

func {{ .Name }}Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
}
`

var httpHeaderTemplate = template.Must(template.New("http_header").Parse(_httpHandlerTemplate))

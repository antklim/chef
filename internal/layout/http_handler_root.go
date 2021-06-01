package layout

import "text/template"

const _httpHeaderRootTemplate = `
package http

import "net/http"

func Muxer(mux *http.ServeMux) {
{{range .Handlers}}
	mux.Handle({{ .Path }}, {{ .Name }}Handler())
{{end}}
}
`

var httpHeaderRootTemplate = template.Must(template.New("http_header_root").Parse(_httpHeaderRootTemplate))

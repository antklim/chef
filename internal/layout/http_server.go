package layout

import "github.com/antklim/chef/internal/layout/template"

var httpServer = newfnode("server.go", withTemplate(template.Get(template.HTTPServer)))

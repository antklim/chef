package layout

import "github.com/antklim/chef/internal/layout/template"

var httpRouter = newfnode("router.go", withTemplate(template.Get(template.HTTPRouter)))

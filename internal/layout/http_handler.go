package layout

import (
	"fmt"

	"github.com/antklim/chef/internal/layout/template"
)

func httpHandler(name string) fnode {
	return newfnode(fmt.Sprintf("%s.go", name), withTemplate(template.Get(template.HTTPHandler)))
}

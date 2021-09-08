package display

import (
	"fmt"
	"io"

	"github.com/antklim/chef/internal/project"
)

const componentListFormat = "%s\t%s\t%s\n"

func ComponentsList(w io.Writer, components []project.Component) error {
	tw.Init(w, minwidth, tabwidth, padding, padchar, flags)
	_, err := fmt.Fprintf(tw, componentListFormat, "NAME", "LOCATION", "DESCRIPTION")
	if err != nil {
		return err
	}

	for _, component := range components {
		_, err := fmt.Fprintf(tw, componentListFormat, component.Name, component.Loc, component.Desc)
		if err != nil {
			return err
		}
	}

	return tw.Flush()
}

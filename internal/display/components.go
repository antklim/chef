package display

import (
	"fmt"
	"io"

	"github.com/antklim/chef/internal/project"
)

const (
	componentsListTitle    = "registered components:"
	componentsListFormat   = "%s\t%s\t%s\n"
	componentsEmptyListMsg = "\tproject does not have registered components"
)

func ComponentsList(w io.Writer, components []project.Component) error {
	ew := &errorWriter{Writer: w}

	fmt.Fprintln(ew, componentsListTitle)

	if len(components) == 0 {
		fmt.Fprintln(ew, componentsEmptyListMsg)
		return ew.err
	}

	tw.Init(ew, minwidth, tabwidth, padding, padchar, flags)

	fmt.Fprintf(tw, componentsListFormat, "NAME", "LOCATION", "DESCRIPTION")

	for _, component := range components {
		fmt.Fprintf(tw, componentsListFormat, component.Name, component.Loc, component.Desc)
	}

	if err := tw.Flush(); err != nil {
		return err
	}

	return ew.err
}

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
	err := componentsList(ew, components)
	if ew.err != nil {
		return ew.err
	}
	return err
}

func componentsList(w io.Writer, components []project.Component) error {
	fmt.Fprintln(w, componentsListTitle)

	if len(components) == 0 {
		fmt.Fprintln(w, componentsEmptyListMsg)
		return nil
	}

	tw.Init(w, minwidth, tabwidth, padding, padchar, flags)

	fmt.Fprintf(tw, componentsListFormat, "NAME", "LOCATION", "DESCRIPTION")

	for _, component := range components {
		fmt.Fprintf(tw, componentsListFormat, component.Name, component.Loc, component.Desc)
	}

	if err := tw.Flush(); err != nil {
		return err
	}
	return nil
}

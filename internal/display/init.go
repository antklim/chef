package display

import (
	"fmt"
	"io"

	"github.com/antklim/chef/internal/project"
)

// TODO (ref): encapsulate location to project and pass project instance to this
// ProjectInit

// ProjectInit outputs information about inited project.
func ProjectInit(w io.Writer, loc string, components []project.Component) error {
	ew := &errorWriter{Writer: w}

	fmt.Fprintf(ew, "project successfully inited at %s\n\n", loc)

	err := componentsList(ew, components)
	if ew.err != nil {
		return ew.err
	}
	return err
}

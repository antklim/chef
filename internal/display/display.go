package display

import (
	"io"
	"text/tabwriter"
)

// the following are tabwriter init parameters
const (
	minwidth      = 0    // minimal cell width including any padding
	tabwidth      = 8    // width of tab characters (equivalent number of spaces)
	padding       = 0    // padding added to a cell before computing its width
	padchar  byte = '\t' // ASCII char used for padding
	flags         = 0    // formatting control
)

// a tabwriter instance
var tw = new(tabwriter.Writer)

// errorWriter is helper structure that provide wraps io.Writer and handles
// errors occurred during Write.
type errorWriter struct {
	io.Writer
	err error
}

func (e *errorWriter) Write(buf []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}

	var n int
	n, e.err = e.Writer.Write(buf)
	return n, e.err
}

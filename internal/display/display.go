package display

import "text/tabwriter"

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

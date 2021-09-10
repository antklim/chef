package chef

import (
	"io"

	"gopkg.in/yaml.v2"
)

type notation struct {
	Version  string
	Category string
	Server   string
}

// Notation defines chef project notation.
type Notation struct {
	Category string
	Server   string
}

func (n Notation) Write(w io.Writer) error {
	nn := notation{
		Version:  version,
		Category: n.Category,
		Server:   n.Server,
	}

	enc := yaml.NewEncoder(w)
	if err := enc.Encode(nn); err != nil {
		return err
	}

	return enc.Close()
}

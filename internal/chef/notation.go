package chef

import (
	"io"

	"gopkg.in/yaml.v2"
)

// Notation defines chef project notation.
type Notation struct {
	Category string
	Server   string
}

// Write writes notation to provided output.
func (n Notation) Write(w io.Writer) error {
	notation := struct {
		Version  string
		Notation `yaml:",inline"`
	}{
		Version:  version,
		Notation: n,
	}

	enc := yaml.NewEncoder(w)
	if err := enc.Encode(notation); err != nil {
		return err
	}

	return enc.Close()
}

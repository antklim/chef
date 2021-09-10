package chef

import (
	"io"

	"gopkg.in/yaml.v2"
)

type notation struct {
	Version  string
	Notation `yaml:",inline"`
}

// Notation defines chef project notation.
type Notation struct {
	Category string
	Server   string
}

// Write writes notation to provided output.
func (n Notation) Write(w io.Writer) error {
	nttn := notation{
		Version:  version,
		Notation: n,
	}

	enc := yaml.NewEncoder(w)
	if err := enc.Encode(nttn); err != nil {
		return err
	}

	return enc.Close()
}

// ReadNotation reads notation from provided source.
func ReadNotation(r io.Reader) (Notation, error) {
	dec := yaml.NewDecoder(r)
	var nttn notation
	err := dec.Decode(&nttn)
	return nttn.Notation, err
}

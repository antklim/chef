package chef

import (
	"io"

	"gopkg.in/yaml.v2"
)

// DefaultNotationName is a default file name to store notation.
const DefaultNotationFileName = ".chef.yml"

// Project describes a project details, like name, language, etc.
type Project struct {
	Name        string
	Description string
	Language    string
}

// Notation defines chef project notation.
type Notation struct {
	Version string
	Project Project
}

// Write writes notation to provided output.
func (n Notation) Write(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	if err := enc.Encode(n); err != nil {
		return err
	}

	return enc.Close()
}

// ReadNotation reads notation from provided source.
func ReadNotation(r io.Reader) (Notation, error) {
	dec := yaml.NewDecoder(r)
	var n Notation
	err := dec.Decode(&n)
	return n, err
}

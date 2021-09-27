package chef

import (
	"io"

	"gopkg.in/yaml.v2"
)

// TODO: add layout information to notation

// TODO: split notation into parts - generic information like category and
// language specific information like module name.

// TODO: add notation builders specializing on language/tech: NewGoNotation ...

// TODO: add language information to notation

// DefaultNotationName is a default file name to store notation.
const DefaultNotationFileName = ".chef.yml"

type notation struct {
	Version  string
	Notation `yaml:",inline"`
}

// Notation defines chef project notation.
type Notation struct {
	Category string
	Server   string `yaml:",omitempty"`
	Module   string `yaml:",omitempty"` // Go module name
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

package cli

type Project interface {
	Bootstrap() error
	Add(string, string) error
	Location() (string, error)
	Name() string
}

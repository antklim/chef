package cli

type Project interface {
	Init() error
	Add(string, string) error
	Location() (string, error)
	Name() string
}

package cli

type Project interface {
	Init() error
	Employ(string, string) error
	Location() (string, error)
	Name() string
}

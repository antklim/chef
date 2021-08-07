package cli

type Project interface {
	Init() error
	EmployComponent(string, string) error
	Location() string
	Name() string
}

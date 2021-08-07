package cli

type Project interface {
	Init() error
	Build() (string, error)
	EmployComponent(string, string) error
}

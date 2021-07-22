package cli

type Project interface {
	Bootstrap() error
	Location() (string, error)
	Name() string
}

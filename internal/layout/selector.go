package layout

import "errors"

type Category string

func (c Category) isUnknown() bool {
	return c != CategoryService
}

const (
	// CategoryService represents service category of a layout.
	CategoryService Category = "srv"
)

var (
	errUnknownCategory = errors.New("unknown layout category")
	errUnknownServer   = errors.New("unknown server")
)

type Server string

func (s Server) isUnknown() bool {
	return s != ServerHTTP
}

const (
	// ServerHTTP represents http server
	ServerHTTP Server = "http"
)

func Selector(c Category, s Server) ([]Node, error) {
	if c.isUnknown() {
		return nil, errUnknownCategory
	}
	if s.isUnknown() {
		return nil, errUnknownServer
	}
	return nil, nil
}

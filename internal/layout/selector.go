package layout

type Category string

const (
	// CategoryUnknown represents unknown category of a layout.
	CategoryUnknown Category = "unknown"
	// CategoryService represents service category of a layout.
	CategoryService Category = "srv"
)

type Server string

const (
	// ServerUnknown represents unknown
	ServerUnknown Server = "unknown"
	// ServerHTTP represents http server
	ServerHTTP Server = "http"
)

func Selector(category, server string) ([]Node, error) {
	return nil, nil
}

package layout

// TODO: finish moving const name to structures

const (
	dirHandler  = "handler"
	dirHTTP     = "http"
	dirInternal = "internal" // nolint
	dirPkg      = "pkg"      // nolint
	dirServer   = "server"
)

var (
	dirAdapter = newdnode("adapter")
	dirApp     = newdnode("app")
	// dirCmd     = newdnode("cmd")
	dirProvider = newdnode("provider")
	dirTest     = newdnode("test")
)

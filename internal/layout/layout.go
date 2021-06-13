package layout

type Layout interface {
	Nodes() []Node
	Schema() string
}

var (
	// m is a map from schema to layout.
	m = make(map[string]Layout)
)

// Register registers the layout to the layouts map. l.Schema will be
// used as the schema registered with this layout.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe.
func Register(l Layout) {
	m[l.Schema()] = l
}

// Get returns the layout registered with the given schema.
//
// If no layout is registered with the schema, nil will be returned.
func Get(name string) Layout {
	if l, ok := m[name]; ok {
		return l
	}
	return nil
}

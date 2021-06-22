package layout

type Layout struct {
	nodes  []Node
	schema string
}

// New creates a new layout with schema s and nodes n.
func New(s string, n []Node) Layout {
	return Layout{
		schema: s,
		nodes:  n,
	}
}

func (l Layout) Nodes() []Node {
	return l.nodes
}

func (l Layout) Schema() string {
	return l.schema
}

func (l Layout) Build(loc string) error {
	for _, n := range l.nodes {
		if err := n.Build(loc); err != nil {
			return err
		}
	}
	return nil
}

const (
	// ServiceLayout an abstract service layout name.
	ServiceLayout = "srv"
	// HTTPServiceLayout an http service layout name.
	HTTPServiceLayout = "srv_http"
)

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
// If no layout is registered with the schema, Nil layout will be returned.
func Get(name string) *Layout {
	if l, ok := m[name]; ok {
		return &l
	}
	return nil
}

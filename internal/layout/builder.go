package layout

// TODO: in imports replace chef/... with the project name

// TODO: read layout settings from yaml
// TODO: test/build generated go code

// TODO: use http handler template to add health endpoint (on bootstrap)
// TODO: make adding health endpoint on bootstrap optional

// TODO: support functionality of bring your own templates

// TODO: init project with go.mod

const (
	fperm = 0644
	dperm = 0755
)

func Builder(root, name string, l Layout) error {
	// root is a project root
	// name is a project name
	n := newdnode(name, withSubNodes(l.Nodes()...))
	return buildNode(root, n)
}

func buildNode(loc string, n Node) error {
	return n.Build(loc)
}

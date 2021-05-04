package chef

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// TODO: read layout settings from yaml
// TODO: refactor project structure
// TODO: use go:embed to init projects

type ProjectCategory string

const (
	CategoryApp ProjectCategory = "app"
	CategoryPkg ProjectCategory = "pkg"
)

type ProjectServer string

const (
	ServerHTTP ProjectServer = "http"
	ServerGRPC ProjectServer = "grpc"
)

type layoutDir int

const (
	dirCmd layoutDir = iota + 1
	dirInternal
	dirTest
	dirApp
	dirAdapter
	dirProvider
	dirServer
	dirHTTP
)

var dirName = map[layoutDir]string{
	dirCmd:      "cmd",
	dirInternal: "internal",
	dirTest:     "test",
	dirApp:      "app",
	dirAdapter:  "adapter",
	dirProvider: "provider",
	dirServer:   "server",
	dirHTTP:     "http",
}

type node int

const (
	nodeDir node = iota
	nodeFile
)

type layoutNode struct {
	Name     string
	Type     node
	Children []layoutNode
}

var defaultLayout = []layoutNode{
	{
		Name: dirName[dirCmd],
		Children: []layoutNode{
			{Name: "main.go", Type: nodeFile},
		},
	},
	{
		Name: dirName[dirInternal],
		Children: []layoutNode{
			{Name: dirName[dirApp]},
			{Name: dirName[dirAdapter]},
			{Name: dirName[dirProvider]},
			{
				Name: dirName[dirServer],
				Children: []layoutNode{
					{Name: dirName[dirHTTP]},
				},
			},
		},
	},
	{Name: dirName[dirTest]},
}

func layoutBuilder(root string, node layoutNode) error {
	o := path.Join(root, node.Name) // file system object, either file or directory

	switch node.Type {
	case nodeFile:
		f, err := os.Create(o)
		if err != nil {
			return err
		}
		return f.Chmod(0644) // nolint
	case nodeDir:
		fallthrough
	default:
		if err := os.Mkdir(o, 0755); err != nil {
			return err
		}

		for _, c := range node.Children {
			if err := layoutBuilder(o, c); err != nil {
				return err
			}
		}
	}

	return nil
}

// Project manager.
type Project struct {
	Name     string
	Root     string
	Category ProjectCategory
	Server   ProjectServer
}

func defaultProject(name string) Project {
	return Project{
		Name:     name,
		Category: CategoryApp,
		Server:   ServerHTTP,
	}
}

// New project.
func New(name string, options ...Option) Project {
	name = strings.TrimSpace(name)
	p := defaultProject(name)
	for _, opt := range options {
		opt.apply(&p)
	}
	return p
}

// TODO: implement option validation
func (p Project) Validate() error {
	if p.Name == "" {
		return errors.New("project name required: empty name provided")
	}

	root, err := p.root()
	if err != nil {
		return err
	}

	fi, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", root)
	}

	fi, _ = os.Stat(path.Join(root, p.Name))
	if fi != nil {
		return fmt.Errorf("file or directory %s already exists", p.Name)
	}

	return nil
}

// Init initializes the project layout.
// TODO: make init a package level function with the default project.Init call.
func (p Project) Init() error {
	rl := layoutNode{
		Name:     p.Name,
		Children: defaultLayout,
	}

	root, err := p.root()
	if err != nil {
		return err
	}

	if err := layoutBuilder(root, rl); err != nil {
		return err
	}

	return nil
}

func (p Project) root() (root string, err error) {
	root = p.Root
	if root == "" {
		root, err = os.Getwd()
	}
	return
}

type Option interface {
	apply(*Project)
}

type funcOption struct {
	f func(*Project)
}

func (fo *funcOption) apply(p *Project) {
	fo.f(p)
}

func newFuncOption(f func(*Project)) *funcOption {
	return &funcOption{f}
}

func WithRoot(r string) Option {
	return newFuncOption(func(p *Project) {
		p.Root = strings.TrimSpace(r)
	})
}

func WithCategory(c ProjectCategory) Option {
	return newFuncOption(func(p *Project) {
		p.Category = c
	})
}

func WithServer(s ProjectServer) Option {
	return newFuncOption(func(p *Project) {
		p.Server = s
	})
}

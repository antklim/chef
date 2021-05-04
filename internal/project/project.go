package project

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/antklim/chef/internal/layout"
)

// TODO: read layout settings from yaml
// TODO: use go:embed to init projects

type Category string

const (
	CategoryApp Category = "app"
	CategoryPkg Category = "pkg"
)

type Server string

const (
	ServerHTTP Server = "http"
	ServerGRPC Server = "grpc"
)

// Project manager.
type Project struct {
	Name     string
	Root     string
	Category Category
	Server   Server
}

func defaultProject(name string) Project {
	return Project{
		Name:     name,
		Category: CategoryApp,
		Server:   ServerHTTP,
	}
}

// New project.
func New(name string, opts ...Option) Project {
	name = strings.TrimSpace(name)
	p := defaultProject(name)
	for _, o := range opts {
		o.apply(&p)
	}
	return p
}

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
func (p Project) Init() error {
	rl := layout.Node{
		Name:     p.Name,
		Children: layout.Default,
	}

	root, err := p.root()
	if err != nil {
		return err
	}

	if err := layout.Builder(root, rl); err != nil {
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

func WithCategory(c Category) Option {
	return newFuncOption(func(p *Project) {
		p.Category = c
	})
}

func WithServer(s Server) Option {
	return newFuncOption(func(p *Project) {
		p.Server = s
	})
}

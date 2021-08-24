package project

import (
	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/layout/node"
	"github.com/antklim/chef/internal/project/template"
)

const (
	dirAdapter  = "adapter"
	dirApp      = "app"
	dirHandler  = "handler"
	dirHTTP     = "http"
	dirServer   = "server"
	dirProvider = "provider"
	dirTest     = "test"
)

type layoutMaker interface {
	makeLayout() *layout.Layout
}

func layoutFactory(category, server string) layoutMaker {
	if category == categoryService && server == serverHTTP {
		return httpServiceLayout{}
	}
	if category == categoryService && server == serverNone {
		return serviceLayout{}
	}
	return nil
}

type serviceLayout struct{}

func (serviceLayout) makeLayout() *layout.Layout {
	nodes := []node.Node{
		node.NewDnode(dirAdapter),
		node.NewDnode(dirApp),
		node.NewDnode(dirHandler),
		node.NewDnode(dirProvider),
		node.NewDnode(dirServer),
		node.NewDnode(dirTest),
	}
	return layout.New(nodes...)
}

type httpServiceLayout struct{}

func (httpServiceLayout) makeLayout() *layout.Layout {
	httpRouter := node.NewFnode("router.go", node.WithTemplate(template.Get(template.HTTPRouter)))
	httpHandlerNode := node.NewDnode(dirHTTP, node.WithSubNodes(httpRouter))
	httpServer := node.NewFnode("server.go", node.WithTemplate(template.Get(template.HTTPServer)))
	httpServerNode := node.NewDnode(dirHTTP, node.WithSubNodes(httpServer))
	httpSrvMain := node.NewFnode("main.go", node.WithTemplate(template.Get(template.HTTPService)))

	nodes := []node.Node{
		node.NewDnode(dirAdapter),
		node.NewDnode(dirApp),
		node.NewDnode(dirHandler, node.WithSubNodes(httpHandlerNode)),
		node.NewDnode(dirProvider),
		node.NewDnode(dirServer, node.WithSubNodes(httpServerNode)),
		node.NewDnode(dirTest),
		httpSrvMain,
	}
	return layout.New(nodes...)
}

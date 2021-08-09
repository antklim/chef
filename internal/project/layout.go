package project

import (
	"github.com/antklim/chef/internal/layout"
	"github.com/antklim/chef/internal/project/template"
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
	nodes := []layout.Node{
		layout.NewDnode(dirAdapter),
		layout.NewDnode(dirApp),
		layout.NewDnode(dirHandler),
		layout.NewDnode(dirProvider),
		layout.NewDnode(dirServer),
		layout.NewDnode(dirTest),
	}
	l := layout.New(ServiceLayout, nodes...)
	return &l
}

type httpServiceLayout struct{}

func (httpServiceLayout) makeLayout() *layout.Layout {
	httpRouter := layout.NewFnode("router.go", layout.WithTemplate(template.Get(template.HTTPRouter)))
	httpHandlerNode := layout.NewDnode(dirHTTP, layout.WithSubNodes(httpRouter))
	httpServer := layout.NewFnode("server.go", layout.WithTemplate(template.Get(template.HTTPServer)))
	httpServerNode := layout.NewDnode(dirHTTP, layout.WithSubNodes(httpServer))
	httpSrvMain := layout.NewFnode("main.go", layout.WithTemplate(template.Get(template.HTTPService)))

	nodes := []layout.Node{
		layout.NewDnode(dirAdapter),
		layout.NewDnode(dirApp),
		layout.NewDnode(dirHandler, layout.WithSubNodes(httpHandlerNode)),
		layout.NewDnode(dirProvider),
		layout.NewDnode(dirServer, layout.WithSubNodes(httpServerNode)),
		layout.NewDnode(dirTest),
		httpSrvMain,
	}
	l := layout.New(HTTPServiceLayout, nodes...)
	return &l
}

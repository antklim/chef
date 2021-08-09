package project

import "fmt"

type IProject interface {
	Init() error
	Employ(string, string) error // employ component
	Name() string
}

type projectMaker func() IProject

var projectMakers = map[string]projectMaker{
	// "svc":      makeSvcProject,
	// "http_svc": makeHttpSvcProject,
}

// type svcProject struct {
// }

// type httpSvcProject struct {
// }

// TODO: it could be layout factory

func MakeProject(category string) (IProject, error) {
	_, ok := projectMakers[category]
	if !ok {
		return nil, fmt.Errorf("unknown project category %q", category)
	}
	return nil, nil
}

// func makeSvcProject() svcProject {
// 	return svcProject{}
// }

// func makeHttpSvcProject() httpSvcProject {
// 	return httpSvcProject{}
// }

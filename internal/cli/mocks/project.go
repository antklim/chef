package mocks

import "github.com/antklim/chef/internal/project"

type Project struct {
	e error
}

func (p Project) Add(c project.Component) error {
	return p.e
}

func (p Project) Bootstrap() error {
	return p.e
}

func (p Project) Location() (string, error) {
	return "", nil
}

func (p Project) Name() string {
	return "ProjectMock"
}

func FailedProject(err error) Project {
	return Project{err}
}

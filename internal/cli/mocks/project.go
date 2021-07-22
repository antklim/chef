package mocks

type Project struct {
	e error
}

func (p Project) Add() error {
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

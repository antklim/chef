package mocks

type Project struct {
	e error
}

func (p Project) EmployComponent(component, name string) error {
	return p.e
}

func (p Project) Init() error {
	return p.e
}

func (p Project) Location() string {
	return ""
}

func (p Project) Name() string {
	return "ProjectMock"
}

func FailedProject(err error) Project {
	return Project{err}
}

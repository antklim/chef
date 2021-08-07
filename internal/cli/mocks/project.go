package mocks

type Project struct {
	initErr  error
	buildErr error
	ecErr    error
	loc      string
}

func (p Project) Init() error {
	return p.initErr
}

func (p Project) Build() (string, error) {
	return p.loc, p.buildErr
}

func (p Project) EmployComponent(component, name string) error {
	return p.ecErr
}

func FailedInit(err error) Project {
	return Project{initErr: err}
}

func FailedBuild(err error) Project {
	return Project{buildErr: err}
}

func FailedEmployComponent(err error) Project {
	return Project{ecErr: err}
}

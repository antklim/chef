package cli

type projMock struct {
	initErr  error
	buildErr error
	ecErr    error
	loc      string
}

func (p projMock) Init() error {
	return p.initErr
}

func (p projMock) Build() (string, error) {
	return p.loc, p.buildErr
}

func (p projMock) EmployComponent(component, name string) error {
	return p.ecErr
}

func FailedInit(err error) Project {
	return projMock{initErr: err}
}

func FailedBuild(err error) Project {
	return projMock{buildErr: err}
}

func FailedEmployComponent(err error) Project {
	return projMock{ecErr: err}
}

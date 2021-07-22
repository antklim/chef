package mocks

type Project struct {
	bootError error
}

func (p Project) Bootstrap() error {
	return p.bootError
}

func (p Project) Location() (string, error) {
	return "", nil
}

func (p Project) Name() string {
	return "ProjectMock"
}

func FailedBootProject(err error) Project {
	return Project{err}
}

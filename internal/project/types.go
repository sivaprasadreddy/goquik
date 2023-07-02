package project

type Project struct {
	ProjectName string `survey:"name"`
	ModulePath  string `survey:"modulePath"`
}

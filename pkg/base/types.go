package base

const (
	SingleSourceProject = 1 + iota
	SingleHeaderProject
	ApplicationProject
	LibraryProject
	XantePluginProject
)

const (
	CLanguage = 1 + iota
	JavaLanguage
	PythonLanguage
	GoLanguage
	RustLanguage
)

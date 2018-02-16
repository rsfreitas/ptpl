package project

import (
	"errors"

	"source-template/pkg/base"
	"source-template/pkg/project/application"
	"source-template/pkg/project/header"
	"source-template/pkg/project/library"
	"source-template/pkg/project/source"
	"source-template/pkg/project/xante"
)

//Our project factory holder
var projectFactory = make(map[int]base.ProjectFactory)

func register(projectType int, project base.ProjectFactory) {
	if project == nil {
		//panic
	}

	_, registered := projectFactory[projectType]

	if registered {
		//error
	}

	projectFactory[projectType] = project
}

//loadSupportedProjects register all supported projects ;-)
func loadSupportedProjects() {
	register(base.SingleSourceProject, source.New)
	register(base.SingleHeaderProject, header.New)
	register(base.ApplicationProject, application.New)
	register(base.LibraryProject, library.New)
	register(base.XantePluginProject, xante.New)
}

// Assemble is responsible to initialize our supported project type and
// build the chosen one.
func Assemble(options base.ProjectOptions) (base.Project, error) {
	loadSupportedProjects()
	project, ok := projectFactory[options.ProjectType]

	if !ok {
		return nil, errors.New("Unimplemented project.")
	}

	return project(options)
}

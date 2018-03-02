//
// Copyright (C) 2017 Rodrigo Freitas
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
//
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
		return nil, errors.New("unimplemented project")
	}

	return project(options)
}

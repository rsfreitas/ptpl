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
package application

import (
	"os"

	"source-template/pkg/base"
	"source-template/pkg/project/common"
	"source-template/pkg/templates"
)

type Application struct {
	// Templates
	sources  []base.FileInfo
	headers  []base.FileInfo
	makefile base.FileInfo

	paths   map[string]string
	Package common.Package
	base.ProjectOptions
}

func (a Application) Build() error {
	// create root path and subdirs
	for _, path := range a.paths {
		err := os.MkdirAll(path, 0755)

		if err != nil {
			return err
		}
	}

	// create sources
	for _, f := range a.sources {
		if err := f.Build(a.paths["source"]); err != nil {
			return err
		}
	}

	// create headers
	for _, f := range a.headers {
		if err := f.Build(a.paths["header"]); err != nil {
			return err
		}
	}

	// create CMakeLists.txt
	if err := a.makefile.Build(a.paths["makefile"]); err != nil {
		return err
	}

	// create package
	if a.PackageProject {
		a.Package.Build()
	}

	return nil
}

func createSources(options base.ProjectOptions) []base.FileInfo {
	var files []base.FileInfo
	sources := []string{
		"main",
	}

	for _, s := range sources {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(s, ".c"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewSource(fileOptions),
		})
	}

	return files
}

func createHeaders(options base.ProjectOptions) []base.FileInfo {
	var files []base.FileInfo
	var headers []string

	for _, suffix := range []string{"_def", "_prt", "_struct"} {
		headers = append(headers, options.ProjectName+suffix)
	}

	headers = append(headers, options.ProjectName)

	for _, h := range headers {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(h, ".h"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewHeader(fileOptions, nil),
		})
	}

	return files
}

func New(options base.ProjectOptions) (base.Project, error) {
	paths := base.Dirtree(options)

	application := &Application{
		ProjectOptions: options,
		paths:          paths,
		sources:        createSources(options),
		headers:        createHeaders(options),
		makefile:       common.CreateMakefile(options),
		Package:        common.NewPackage(options, paths),
	}

	return application, nil
}

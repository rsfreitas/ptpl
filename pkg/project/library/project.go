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
package library

import (
	"os"

	"source-template/pkg/base"
	"source-template/pkg/project/common"
	"source-template/pkg/templates"
)

type Library struct {
	sources  []base.FileInfo
	headers  []base.FileInfo
	makefile base.FileInfo
	symbol   base.FileInfo

	paths   map[string]string
	Package common.Package
	base.ProjectOptions
}

func (l Library) Build() error {
	// create root path and subdirs
	for _, path := range l.paths {
		err := os.MkdirAll(path, 0755)

		if err != nil {
			return err
		}
	}

	// create sources
	for _, f := range l.sources {
		if err := f.Build(l.paths["source"]); err != nil {
			return err
		}
	}

	// create headers
	for _, f := range l.headers {
		if err := f.Build(l.paths["header"]); err != nil {
			return err
		}
	}

	// create CMakeLists.txt
	if err := l.makefile.Build(l.paths["makefile"]); err != nil {
		return err
	}

	// create symbols file
	if err := l.symbol.Build(l.paths["misc"]); err != nil {
		return err
	}

	// create package
	if l.PackageProject {
		l.Package.Build()
	}

	return nil
}

func createSources(options base.ProjectOptions) ([]base.FileInfo, []string) {
	var files []base.FileInfo
	sources := []string{
		"utils",
		"error",
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

	return files, sources
}

func createHeaders(options base.ProjectOptions, sources []string) []base.FileInfo {
	var files []base.FileInfo
	headers := []string{
		"lib" + options.ProjectName,
		"internal/internal.h",
		"internal/utils.h",
		"internal/error.h",
		"api/utils.h",
		"api/error.h",
	}

	for _, h := range headers {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(h, ".h"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewHeader(fileOptions, sources),
		})
	}

	return files
}

func createSymbol(options base.ProjectOptions) base.FileInfo {
	fileOptions := base.FileOptions{
		Executable:     false,
		HeaderComment:  false,
		ProjectOptions: options,
		Name:           "lib" + options.ProjectName + ".sym",
	}

	return base.FileInfo{
		FileOptions:  fileOptions,
		FileTemplate: templates.NewSymbol(fileOptions),
	}
}

func New(options base.ProjectOptions) (base.Project, error) {
	sources, sourceFilenames := createSources(options)
	paths := base.Dirtree(options)

	return &Library{
		sources:        sources,
		paths:          paths,
		ProjectOptions: options,
		headers:        createHeaders(options, sourceFilenames),
		makefile:       common.CreateMakefile(options),
		symbol:         createSymbol(options),
		Package:        common.NewPackage(options, paths),
	}, nil
}

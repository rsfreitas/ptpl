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
package xante

import (
	"os"

	"source-template/pkg/base"
	"source-template/pkg/project/common"
	"source-template/pkg/templates"
)

type XantePlugin struct {
	sources  []base.FileInfo
	headers  []base.FileInfo
	makefile base.FileInfo
	script   base.FileInfo

	paths   map[string]string
	Package common.Package
	base.ProjectOptions
}

func (x XantePlugin) Build() error {
	// create root path and subdirs
	for _, path := range x.paths {
		err := os.MkdirAll(path, 0755)

		if err != nil {
			return err
		}
	}

	// create sources
	for _, f := range x.sources {
		if err := f.Build(x.paths["source"]); err != nil {
			return err
		}
	}

	// create headers
	for _, f := range x.headers {
		if err := f.Build(x.paths["header"]); err != nil {
			return err
		}
	}

	// create CMakeLists.txt
	if err := x.makefile.Build(x.paths["makefile"]); err != nil {
		return err
	}

	// create application script
	if err := x.script.Build(x.paths["script"]); err != nil {
		return err
	}

	// create package
	if x.PackageProject {
		x.Package.Build()
	}

	return nil
}

func createSources(options base.ProjectOptions) []base.FileInfo {
	var extension string
	var files []base.FileInfo
	sources := []string{
		"plugin",
	}

	if options.Language == base.GoLanguage {
		extension = ".go"
	} else {
		extension = ".c"
	}

	for _, s := range sources {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(s, extension),
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
	headers := []string{
		"plugin.h",
	}

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

func createPluginScript(options base.ProjectOptions) base.FileInfo {
	fileOptions := base.FileOptions{
		Executable:     true,
		HeaderComment:  true,
		ProjectOptions: options,
		Name:           options.ProjectName,
	}

	return base.FileInfo{
		FileOptions:  fileOptions,
		FileTemplate: templates.NewBash(fileOptions),
	}
}

func New(options base.ProjectOptions) (base.Project, error) {
	var headers []base.FileInfo
	paths := base.Dirtree(options)

	// Only C plugins have header files
	if options.Language == base.CLanguage {
		headers = createHeaders(options)
	}

	return &XantePlugin{
		paths:          paths,
		sources:        createSources(options),
		headers:        headers,
		ProjectOptions: options,
		makefile:       common.CreateMakefile(options),
		script:         createPluginScript(options),
		Package:        common.NewPackage(options, paths),
	}, nil
}

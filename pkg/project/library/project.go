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
	"fmt"
	"os"

	"source-template/pkg/base"
	"source-template/pkg/project/common"
	"source-template/pkg/templates"
)

type Library struct {
	sources  []base.FileInfo
	headers  []base.FileInfo
	debian   []base.FileInfo
	makefile base.FileInfo
	symbol   base.FileInfo
	rootPath string
	base.ProjectOptions
}

func (l Library) String() string {
	return fmt.Sprintf("Library project")
}

func createLibraryDirtree(path string, options base.ProjectOptions) error {
	var subdirs []string
	var prefix string

	if options.PackageProject {
		prefix = options.ProjectName
		subdirs = append(subdirs, "pkg_install/misc")
		subdirs = append(subdirs, "pkg_install/debian")
	}

	subdirs = append(subdirs, prefix+"/src")
	subdirs = append(subdirs, prefix+"/include/api")
	subdirs = append(subdirs, prefix+"/include/internal")
	subdirs = append(subdirs, prefix+"/misc")

	for _, dir := range subdirs {
		err := os.MkdirAll(path+"/"+dir, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}

func (l Library) Build() error {
	// create root path and subdirs
	if err := createLibraryDirtree(l.rootPath, l.ProjectOptions); err != nil {
		return err
	}

	// create sources
	for _, f := range l.sources {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create headers
	for _, f := range l.headers {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create CMakeLists.txt
	if err := l.makefile.Build(); err != nil {
		return err
	}

	// create symbols file
	if err := l.symbol.Build(); err != nil {
		return err
	}

	return nil
}

func createSources(options base.ProjectOptions, rootPath string, prefix string) ([]base.FileInfo, []string) {
	var files []base.FileInfo
	sources := []string{
		"utils",
		"error",
	}

	for _, s := range sources {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(rootPath+"/"+prefix+"/src/"+s, ".c"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewSource(fileOptions),
		})
	}

	return files, sources
}

func createHeaders(options base.ProjectOptions, rootPath string, prefix string, sources []string) []base.FileInfo {
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
			Name:           base.AddExtension(rootPath+"/"+prefix+"/include/"+h, ".h"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewHeader(fileOptions, sources),
		})
	}

	return files
}

func createSymbol(options base.ProjectOptions, rootPath string, prefix string) base.FileInfo {
	fileOptions := base.FileOptions{
		Executable:     false,
		HeaderComment:  false,
		ProjectOptions: options,
		Name:           rootPath + "/" + prefix + "/misc/lib" + options.ProjectName + ".sym",
	}

	return base.FileInfo{
		FileOptions:  fileOptions,
		FileTemplate: templates.NewSymbol(fileOptions),
	}
}

func New(options base.ProjectOptions) (base.Project, error) {
	var rootPath string
	var prefix string
	cwd, err := os.Getwd()

	if err != nil {
		return &Library{}, err
	}

	if options.PackageProject {
		prefix = options.ProjectName
		rootPath = cwd + "/package-lib" + options.ProjectName
	} else {
		rootPath = cwd + "/lib" + options.ProjectName
	}

	sources, sourceFilenames := createSources(options, rootPath, prefix)

	return &Library{
		rootPath:       rootPath,
		sources:        sources,
		headers:        createHeaders(options, rootPath, prefix, sourceFilenames),
		debian:         common.CreateDebianScripts(options, rootPath),
		ProjectOptions: options,
		makefile:       common.CreateMakefile(options, rootPath, prefix),
		symbol:         createSymbol(options, rootPath, prefix),
	}, nil
}

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
	debian   []base.FileInfo
	makefile base.FileInfo
	script   base.FileInfo
	rootPath string
	base.ProjectOptions
}

func (x XantePlugin) String() string {
	return ""
}

func createPluginDirtree(path string, options base.ProjectOptions) error {
	var subdirs []string
	var prefix string

	if options.PackageProject {
		prefix = options.ProjectName
		subdirs = append(subdirs, "pkg_install/misc")
		subdirs = append(subdirs, "pkg_install/debian")
	}

	subdirs = append(subdirs, prefix+"/src")
	subdirs = append(subdirs, prefix+"/include")
	subdirs = append(subdirs, prefix+"/jtf")
	subdirs = append(subdirs, prefix+"/script")

	for _, dir := range subdirs {
		err := os.MkdirAll(path+"/"+dir, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}

func (x XantePlugin) Build() error {
	if err := createPluginDirtree(x.rootPath, x.ProjectOptions); err != nil {
		return err
	}

	// create sources
	for _, f := range x.sources {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create headers
	for _, f := range x.headers {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create CMakeLists.txt
	if err := x.makefile.Build(); err != nil {
		return err
	}

	// create application script
	if err := x.script.Build(); err != nil {
		return err
	}

	return nil
}

func createSources(options base.ProjectOptions, rootPath string, prefix string) []base.FileInfo {
	var files []base.FileInfo
	sources := []string{
		"plugin",
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

	return files
}

func createHeaders(options base.ProjectOptions, rootPath string, prefix string) []base.FileInfo {
	var files []base.FileInfo
	headers := []string{
		"plugin.h",
	}

	for _, h := range headers {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(rootPath+"/"+prefix+"/include/"+h, ".h"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewHeader(fileOptions, nil),
		})
	}

	return files
}

func createPluginScript(options base.ProjectOptions, rootPath string, prefix string) base.FileInfo {
	fileOptions := base.FileOptions{
		Executable:     true,
		HeaderComment:  true,
		ProjectOptions: options,
		Name:           rootPath + "/" + prefix + "/script/" + options.ProjectName,
	}

	return base.FileInfo{
		FileOptions:  fileOptions,
		FileTemplate: templates.NewBash(fileOptions),
	}
}

func New(options base.ProjectOptions) (base.Project, error) {
	var rootPath string
	var prefix string
	cwd, err := os.Getwd()

	if err != nil {
		return &XantePlugin{}, err
		return &XantePlugin{}, nil
	}

	if options.PackageProject {
		prefix = options.ProjectName
		rootPath = cwd + "/package-lib" + options.ProjectName
	} else {
		rootPath = cwd + "/" + options.ProjectName
	}

	return &XantePlugin{
		rootPath:       rootPath,
		sources:        createSources(options, rootPath, prefix),
		headers:        createHeaders(options, rootPath, prefix),
		debian:         common.CreateDebianScripts(options, rootPath),
		ProjectOptions: options,
		makefile:       common.CreateMakefile(options, rootPath, prefix),
		script:         createPluginScript(options, rootPath, prefix),
	}, nil
}
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
package common

import (
	"strings"

	"source-template/pkg/base"
	"source-template/pkg/templates"
)

type Package struct {
	debian  []base.FileInfo
	service base.FileInfo
	builder base.FileInfo
	options base.ProjectOptions
	paths   map[string]string
}

// Build builds all required and necessary package contents and structure.
func (p *Package) Build() error {
	// create debian scripts
	for _, f := range p.debian {
		if err := f.Build(p.paths["debian"]); err != nil {
			return err
		}
	}

	if err := p.service.Build(p.paths["misc"]); err != nil {
		return err
	}

	if err := p.builder.Build(p.paths["package"]); err != nil {
		return err
	}

	return nil
}

func createDebianScripts(options base.ProjectOptions) []base.FileInfo {
	var files []base.FileInfo
	scripts := []string{
		"preinst",
		"prerm",
		"postinst",
		"postrm",
	}

	// If we're not a package
	if !options.PackageProject {
		return files
	}

	for _, s := range scripts {
		fileOptions := base.FileOptions{
			Executable:     true,
			HeaderComment:  true,
			ProjectOptions: options,
			Name:           s,
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewBash(fileOptions),
		})
	}

	return files
}

func createSystemdService(options base.ProjectOptions) base.FileInfo {
	fileOptions := base.FileOptions{
		Executable:     false,
		HeaderComment:  false,
		ProjectOptions: options,
		Name:           strings.ToLower(options.ProjectName) + ".service",
	}

	return base.FileInfo{
		FileOptions:  fileOptions,
		FileTemplate: templates.NewText(fileOptions),
	}
}

func createBuildScript(options base.ProjectOptions) base.FileInfo {
	fileOptions := base.FileOptions{
		Executable:     true,
		HeaderComment:  true,
		ProjectOptions: options,
		Name:           "build-package.sh",
	}

	return base.FileInfo{
		FileOptions:  fileOptions,
		FileTemplate: templates.NewBash(fileOptions),
	}
}

func NewPackage(options base.ProjectOptions, paths map[string]string) Package {
	return Package{
		options: options,
		paths:   paths,
		debian:  createDebianScripts(options),
		service: createSystemdService(options),
		builder: createBuildScript(options),
	}
}

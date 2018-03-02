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
	"source-template/pkg/base"
	"source-template/pkg/templates"
)

func CreateDebianScripts(options base.ProjectOptions, rootPath string) []base.FileInfo {
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
			Name:           rootPath + "/pkg_install/debian/" + s,
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewBash(fileOptions),
		})
	}

	return files
}

func CreateMakefile(options base.ProjectOptions, rootPath string, prefix string) base.FileInfo {
	fileOptions := base.FileOptions{
		Executable:     false,
		HeaderComment:  false,
		ProjectOptions: options,
		Name:           rootPath + "/" + prefix + "/CMakeLists.txt",
	}

	return base.FileInfo{
		FileOptions:  fileOptions,
		FileTemplate: templates.NewMakefile(fileOptions),
	}
}

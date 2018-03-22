// The project's interface
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
package base

import (
	"os"
)

type ProjectOptions struct {
	PackageProject         bool
	ProjectName            string
	AuthorName             string
	Language               int
	ProjectType            int
	LibcollectionsFeatures bool
}

type Project interface {
	// Build is where all the magic must happen and the template project must
	// be created.
	Build() error
}

// Also, every supported project must have at least a function with the following
// signature:
type ProjectFactory func(ProjectOptions) (Project, error)

// Dirtree fills a map with all needed project sub-directories.
func Dirtree(options ProjectOptions) map[string]string {
	var prefix string
	var rootPath string
	dirtree := make(map[string]string)

	cwd, err := os.Getwd()

	if err != nil {
		return nil
	}

	if options.PackageProject {
		prefix = options.ProjectName
		rootPath = cwd + "/package-" + options.ProjectName
		dirtree["package"] = rootPath + "/pkg_install"
		dirtree["debian"] = rootPath + "/pkg_install/debian"
		dirtree["misc"] = rootPath + "/pkg_install/misc"
	} else {
		rootPath = cwd + "/" + options.ProjectName
	}

	dirtree["source"] = rootPath + "/" + prefix + "/src"

	if options.Language == CLanguage {
		dirtree["header"] = rootPath + "/" + prefix + "/include"
	}

	if options.ProjectType == XantePluginProject {
		dirtree["script"] = rootPath + "/" + prefix + "/script"
		dirtree["jtf"] = rootPath + "/" + prefix + "/jtf"
		dirtree["makefile"] = dirtree["source"]
	} else {
		dirtree["makefile"] = rootPath + "/" + prefix
	}

	if options.ProjectType == LibraryProject {
		dirtree["api-header"] = rootPath + "/" + prefix + "/include/api"
		dirtree["internal-header"] = rootPath + "/" + prefix + "/include/internal"
		dirtree["misc"] = rootPath + "/" + prefix + "/misc"
	}

	return dirtree
}

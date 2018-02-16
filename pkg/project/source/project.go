// The single source project type implementation.
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
package source

import (
	"fmt"

	"source-template/pkg/base"
	"source-template/pkg/templates"
)

type SingleSource struct {
	file base.FileInfo
	base.ProjectOptions
}

func (s SingleSource) String() string {
	return fmt.Sprintf("Project type: single source\nFilename: %s\n", s.ProjectName)
}

func (s SingleSource) Build() error {
	return s.file.Build()
}

func New(options base.ProjectOptions) (base.Project, error) {
	fileOptions := base.FileOptions{
		Name:           base.AddExtension(options.ProjectName, ".c"),
		HeaderComment:  true,
		ProjectOptions: options,
	}

	return &SingleSource{
		file: base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewSource(fileOptions),
		},
		ProjectOptions: options,
	}, nil
}

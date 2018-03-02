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
	"source-template/pkg/base"
)

type XantePlugin struct {
	sources  []base.FileInfo
	headers  []base.FileInfo
	rootPath string
	base.ProjectOptions
}

func (l XantePlugin) String() string {
	return ""
}

func (l XantePlugin) Build() error {
	// create root path and subdirs
	// create sources
	// create headers
	// create Makefile (future CMakeLists.txt)
	// create application script

	return nil
}

func New(options base.ProjectOptions) (base.Project, error) {
	return &XantePlugin{}, nil
}

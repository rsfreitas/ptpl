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
package templates

import (
	"os"
	"text/template"

	"source-template/pkg/base"
)

type SymbolFile struct {
	base.FileOptions
	ContentData
}

const content = `LIB{{.ProjectNameUpper}}_0.1 {
	global:
		*;
	local:
		*;
};
`

func (s SymbolFile) Header(file *os.File) {
}

func (s SymbolFile) HeaderComment(file *os.File) {
}

func (s SymbolFile) Footer(file *os.File) {
}

func (s SymbolFile) Content(file *os.File) {
	tmpTpl := template.New("symbol")
	tpl, err := tmpTpl.Parse(content)

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func NewSymbol(options base.FileOptions) base.FileTemplate {
	return &SymbolFile{
		FileOptions: options,
		ContentData: GetContentData(options),
	}
}

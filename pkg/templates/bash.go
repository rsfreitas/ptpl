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
	"fmt"
	"os"

	"source-template/pkg/base"
)

type BashFile struct {
	base.FileOptions
	ContentData
}

func (s BashFile) Header(file *os.File) {
	// if we're creating a project, probably will have an include directive here
}

func (s BashFile) HeaderComment(file *os.File) {
	tpl, err := BashSourceHeader()

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func (s BashFile) Footer(file *os.File) {
	file.WriteString("\nexit 0\n")
}

func (s BashFile) Content(file *os.File) {
	fmt.Println("Single source content")
	// here we try to guess what will be the file content by its name (basename)
}

func NewBash(options base.FileOptions) base.FileTemplate {
	return &BashFile{
		FileOptions: options,
		ContentData: GetContentData(options),
	}
}

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

const serviceContent = `
[Unit]
Description=
After=

[Service]
Type=simple
User=root
WorkingDirectory=
ExecStart=
Restart=always
RestartSec=1

[Install]
WantedBy=multi-user.target
`

type TextFile struct {
	content string
	base.FileOptions
	ContentData
}

func (s TextFile) Header(file *os.File) {
}

func (s TextFile) HeaderComment(file *os.File) {
}

func (s TextFile) Footer(file *os.File) {
}

func (s TextFile) Content(file *os.File) {
	tmpTpl := template.New("script")
	tpl, err := tmpTpl.Parse(s.content)

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func NewText(options base.FileOptions) base.FileTemplate {
	var content string
	_, extension := extractFilename(options.Name, options.ProjectType)

	if options.PackageProject {
		if extension == ".service" {
			content = serviceContent
		}
	}

	return &TextFile{
		FileOptions: options,
		content:     content,
		ContentData: GetContentData(options),
	}
}

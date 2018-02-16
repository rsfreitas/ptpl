// The C file source type implementation.
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
	"path/filepath"

	"source-template/pkg/base"
)

type SourceFile struct {
	filename string
	content  string
	base.FileOptions
}

func (s SourceFile) Header(file *os.File) {
	// if we're creating a project, probably will have an include directive here
	cnt := fmt.Sprintf("\n#include \"%[1]s.h\"\n", s.filename)
	file.WriteString(cnt)
}

func (s SourceFile) HeaderComment(file *os.File) {
	file.WriteString(`
/*
 * Description:
 *
 * Author:
 * Created at:
 * Project:
 *
 * Copyright (C) 2017 Author Name All rights reserved.
 */
`)
}

func (s SourceFile) Footer(file *os.File) {
	//nothing here
}

func (s SourceFile) Content(file *os.File) {
	if s.content != "" {
		file.WriteString(s.content)
	}
}

func mainContent(projectName string) string {
	return `
int main(int argc, char **argv)
{
	const char *opt = "hv\0";
	int option;

	do {
		option = getopt(argc, argv, opt);

		switch (option) {
			case 'h':
				return 1;

			case 'v':
				return 1;

			case '?':
				return -1;
		}
	} while (option != -1);

	return 0;
}

`
}

func NewSource(options base.FileOptions) base.FileTemplate {
	var content string

	bname := filepath.Base(options.Name)
	extension := filepath.Ext(bname)
	bname = bname[0 : len(bname)-len(extension)]

	// here we build what will be the file content based on its name (basename)
	if options.ProjectType == base.ApplicationProject {
		if bname == "main" {
			content = mainContent(options.ProjectName)
		}
	}

	return &SourceFile{
		FileOptions: options,
		filename:    bname,
		content:     content,
	}
}

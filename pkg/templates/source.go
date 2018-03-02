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
	"text/template"

	"source-template/pkg/base"
)

type SourceFile struct {
	filename string
	content  string
	options  base.FileOptions
	ContentData
}

func (s SourceFile) Header(file *os.File) {
	var cnt string

	// if we're creating a project, probably will have an include directive here
	if s.options.ProjectType == base.LibraryProject {
		cnt = fmt.Sprintf("\n#include \"lib%[1]s.h\"\n", s.options.ProjectName)
	} else {
		cnt = fmt.Sprintf("\n#include \"%[1]s.h\"\n", s.options.ProjectName)
	}

	file.WriteString(cnt)
}

func (s SourceFile) HeaderComment(file *os.File) {
	tpl, err := CSourceHeader()

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func (s SourceFile) Footer(file *os.File) {
	//nothing here
}

func (s SourceFile) Content(file *os.File) {
	tmpTpl := template.New("source")
	tpl, err := tmpTpl.Parse(s.content)

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func mainContent() string {
	return `
static void usage(void)
{
    printf("Usage: %s [OPTIONS]\n", APP_NAME);
    printf("A brief description.\n\n");
    printf("Options:\n\n");
    printf("  -h, --help                 Shows this help screen.\n");
    printf("  -v, --version              Shows current jerminus version.\n");
    printf("\n");
}

static void version(void)
{
    printf("%s - Version %d.%d.%d %s\n", APP_NAME, MAJOR_VERSION, MINOR_VERSION,
           RELEASE, (BETA == true) ? "beta" : "");
}

int main(int argc, char **argv)
{
	const char *opt = "hv\0";
	int option;

	do {
		option = getopt(argc, argv, opt);

		switch (option) {
			case 'h':
				usage();
				return 1;

			case 'v':
				version();
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
	bname := extractFilename(options.Name, options.ProjectType)
	contentData := GetContentData(options)

	// here we build what will be the file content based on its name (basename)
	if bname == "main" {
		content = mainContent()
	} else if bname == "error" {
		content = errorContent(Source, options)
	}

	return &SourceFile{
		options:     options,
		filename:    bname,
		content:     content,
		ContentData: contentData,
	}
}

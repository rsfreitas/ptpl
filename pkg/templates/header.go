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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"source-template/pkg/base"
)

type HeaderFile struct {
	filename   string // The header file basename.
	content    string // A custom header content.
	headerPath string
	base.FileOptions
	ContentData
}

func (s HeaderFile) Header(file *os.File) {
	upper := strings.Replace(s.filename, "-", "_", -1)

	if s.ProjectType == base.LibraryProject {
		upper = fmt.Sprintf("LIB%s_%s_%s", s.ProjectName, s.headerPath, upper)
	}

	cnt := fmt.Sprintf("\n#ifndef _%[1]s_H\n#define _%[1]s_H\n", strings.ToUpper(upper))
	file.WriteString(cnt)
}

func (s HeaderFile) HeaderComment(file *os.File) {
	tpl, err := CSourceHeader()

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func (s HeaderFile) Footer(file *os.File) {
	file.WriteString("\n#endif\n")
}

func (s HeaderFile) Content(file *os.File) {
	tpl := template.New("header")

	tpl, err := tpl.Parse(s.content)

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

// applicationMainHeaderContent builds the content (body) of the main header file
// of a project.
func applicationMainHeaderContent(name string, projectType int) string {
	var cnt string

	if projectType == base.LibraryProject {
		cnt = fmt.Sprintf(`
{{.LibcollectionsInclude}}

#ifdef LIB%[1]s_COMPILE
# define MAJOR_VERSION		0
# define MINOR_VERSION		1
# define RELEASE			1

# include "internal/internal.h"
#endif

{{.ProjectIncludeFiles}}`, strings.ToUpper(name))
	} else if projectType == base.ApplicationProject {
		cnt = fmt.Sprintf(`
/* Standard library headers */
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdbool.h>

/* External library headers */
{{.LibcollectionsInclude}}

/* Internal headers */
#include "%[1]s_def.h"
#include "%[1]s_struct.h"
#include "%[1]s_prt.h"
`, name)
	}

	return cnt
}

const internalLibraryHeaderContent = `
/*
 * An internal representation of a public function. It does not affect the code
 * or the function visibility. Its objective is only to let clear what is and
 * what is not been exported from library by looking at the code.
 *
 * Every exported function must have this at the beginning of its declaration.
 * Example:
 *
 * __PUB_API__ const char *function(int arg)
 * {
 *      // Body
 * }
 */
#define __PUB_API__

/* Internal library API */
{{.ProjectIncludeFiles}}`

const applicationDefines = `
#define MAJOR_VERSION			0
#define MINOR_VERSION			1
#define RELEASE					1
#define BETA					true

#define APP_NAME				"{{.ProjectName}}"
`

const pluginHeaderContent = `
/* External libraries */
#include <collections.h>
#include <libxante.h>
`

func projectIncludeFiles(sourceFilenames []string, includePath string) string {
	var s bytes.Buffer

	for _, source := range sourceFilenames {
		s.WriteString(fmt.Sprintf("#include \"%s/%s.h\"\n", includePath, source))
	}

	return s.String()
}

// NewHeader creates C header file template. It must receive the desired file
// options, containing informations about it. It also receives a list of source
// file names to be used in special cases, such as building the include files
// preprocessor of a library.
func NewHeader(options base.FileOptions, sources []string) base.FileTemplate {
	var content, dir, includePath string
	bname := extractFilename(options.Name, options.ProjectType)

	if bname == options.ProjectName {
		content = applicationMainHeaderContent(bname, options.ProjectType)
		includePath = "api"
	} else if bname == "internal" {
		content = internalLibraryHeaderContent
		includePath = "internal"
	} else if bname == "error" {
		dir = filepath.Base(filepath.Dir(options.Name))
		flags := ContentType(0)

		if dir == "internal" {
			flags = InternalHeader
		}

		content = errorContent(flags, options)
	} else if strings.Contains(bname, "_def") {
		content = applicationDefines
	} else if bname == "plugin" {
		content = pluginHeaderContent
	}

	contentData := GetContentData(options)

	if options.ProjectType == base.LibraryProject {
		contentData.ProjectIncludeFiles = projectIncludeFiles(sources, includePath)
	}

	if options.LibcollectionsFeatures {
		if options.ProjectType == base.LibraryProject {
			contentData.LibcollectionsInclude = `#ifndef _COLLECTIONS_H
# include <collections.h>
#endif`
		} else {
			contentData.LibcollectionsInclude = "#include <collections.h>"
		}
	}

	return &HeaderFile{
		FileOptions: options,
		content:     content,
		filename:    bname,
		headerPath:  dir,
		ContentData: contentData,
	}
}

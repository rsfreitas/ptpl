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
	upper := s.filename

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

// applicationMainHeaderContent builds the content (body) of the main header file.
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
	} else {
		cnt = fmt.Sprintf("\n/* Standard library headers */\n"+
			"#include <stdlib.h>\n"+
			"#include <unistd.h>\n"+
			"\n/* External library headers */\n"+
			"\n/* Internal headers */\n"+
			"#include \"%[1]s_def.h\"\n"+
			"#include \"%[1]s_struct.h\"\n"+
			"#include \"%[1]s_prt.h\"\n", name)
	}

	return cnt
}

func internalLibraryHeaderContent() string {
	return `
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
}

func projectIncludeFiles(sourceFilenames []string, includePath string) string {
	var s bytes.Buffer

	for _, source := range sourceFilenames {
		s.WriteString(fmt.Sprintf("#include \"%s/%s.h\"\n", includePath, source))
	}

	return s.String()
}

// NewHeader creates C header file template. It must receive the wanted file
// options, containing informations about it, and a custom content, as the "body"
// of the new file.
func NewHeader(options base.FileOptions, sources []string) base.FileTemplate {
	var content, dir, includePath string
	bname := extractFilename(options.Name, options.ProjectType)

	if bname == options.ProjectName {
		content = applicationMainHeaderContent(bname, options.ProjectType)
		includePath = "api"
	} else if bname == "internal" {
		content = internalLibraryHeaderContent()
		includePath = "internal"
	} else if bname == "error" {
		dir = filepath.Base(filepath.Dir(options.Name))
		flags := ContentType(0)

		if dir == "internal" {
			flags = InternalHeader
		}

		content = errorContent(flags, options)
	}

	contentData := GetContentData(options)

	if options.ProjectType == base.LibraryProject {
		contentData.ProjectIncludeFiles = projectIncludeFiles(sources, includePath)
	}

	if options.LibcollectionsFeatures {
		contentData.LibcollectionsInclude = `#ifndef _COLLECTIONS_H
# include <collections.h>
#endif`
	}

	return &HeaderFile{
		FileOptions: options,
		content:     content,
		filename:    bname,
		headerPath:  dir,
		ContentData: contentData,
	}
}

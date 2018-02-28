package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"source-template/pkg/base"
)

type HeaderFile struct {
	filename string //The header file basename.
	content  string //A custom header content.
	base.FileOptions
	ContentData
}

func (s HeaderFile) Header(file *os.File) {
	upper := strings.ToUpper(s.filename)

	if s.ProjectType == base.LibraryProject {
		upper = "LIB" + upper
	}

	cnt := fmt.Sprintf("\n#ifndef _%s_H\n#define _%s_H\n", upper, upper)
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
	if s.content != "" {
		file.WriteString(s.content)
	}
}

//applicationMainHeaderContent builds the content (body) of the main header file.
func applicationMainHeaderContent(name string, projectType int) string {
	var cnt string

	if projectType == base.LibraryProject {
		cnt = fmt.Sprintf(`
#ifdef LIB%[1]s_COMPILE
# define MAJOR_VERSION		0
# define MINOR_VERSION		1
# define RELEASE			1

# include "internal/internal.h"
#endif
		`, strings.ToUpper(name))
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

// extractFilename gives only the file name without path and extension.
func extractFilename(filename string, projectType int) string {
	bname := filepath.Base(filename)
	extension := filepath.Ext(bname)
	bname = bname[0 : len(bname)-len(extension)]

	if projectType == base.LibraryProject && strings.Contains(bname, "lib") {
		bname = bname[3:]
	}

	fmt.Println(bname)
	return bname
}

// TODO: Remove content argument
//NewHeader creates C header file template. It must receive the wanted file
//options, containing informations about it, and a custom content, as the "body"
//of the new file.
func NewHeader(options base.FileOptions, content string) base.FileTemplate {
	bname := extractFilename(options.Name, options.ProjectType)

	if bname == options.ProjectName {
		content = applicationMainHeaderContent(bname, options.ProjectType)
	}

	return &HeaderFile{
		FileOptions: options,
		content:     content,
		filename:    bname,
		ContentData: GetContentData(options),
	}
}

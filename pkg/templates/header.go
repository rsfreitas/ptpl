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
}

func (s HeaderFile) Header(file *os.File) {
	upper := strings.ToUpper(s.filename)
	cnt := fmt.Sprintf("\n#ifndef _%s_H\n#define _%s_H\n", upper, upper)
	file.WriteString(cnt)
}

func (s HeaderFile) HeaderComment(file *os.File) {
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

func (s HeaderFile) Footer(file *os.File) {
	file.WriteString("\n#endif\n")
}

func (s HeaderFile) Content(file *os.File) {
	if s.content != "" {
		file.WriteString(s.content)
	}
}

//applicationMainHeaderContent builds the content (body) of the main header file.
func applicationMainHeaderContent(name string) string {
	cnt := fmt.Sprintf("\n/* Standard library headers */\n"+
		"#include <stdlib.h>\n"+
		"#include <unistd.h>\n"+
		"\n/* External library headers */\n"+
		"\n/* Internal headers */\n"+
		"#include \"%[1]s_def.h\"\n"+
		"#include \"%[1]s_struct.h\"\n"+
		"#include \"%[1]s_prt.h\"\n", name)

	return cnt
}

//NewHeader creates C header file template. It must receive the wanted file
//options, containing informations about it, and a custom content, as the "body"
//of the new file.
func NewHeader(options base.FileOptions, content string) base.FileTemplate {
	//saves only the filename, without path and extension
	bname := filepath.Base(options.Name)
	extension := filepath.Ext(bname)
	bname = bname[0 : len(bname)-len(extension)]

	if bname == options.ProjectName {
		content = applicationMainHeaderContent(bname)
	}

	return &HeaderFile{
		FileOptions: options,
		content:     content,
		filename:    bname,
	}
}

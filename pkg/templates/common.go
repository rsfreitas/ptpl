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
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"source-template/pkg/base"
)

type ContentType int

const (
	Source ContentType = 1 << iota
	InternalHeader
)

const bashHeaderContent = `
#
# Description:
#
# Author: {{.Author}}
# Created at: {{.Date}}
# Project: {{.ProjectName}}
#
# Copyright (C) {{.Year}} {{.Author}} All rights reserved.
#
`

const headerContent = `
/*
 * Description:
 *
 * Author: {{.Author}}
 * Created at: {{.Date}}
 * Project: {{.ProjectName}}
 *
 * Copyright (C) {{.Year}} {{.Author}} All rights reserved.
 */
`

// ContentData must be used to replace variables inside template strings.
type ContentData struct {
	ProjectName           string
	Author                string
	Date                  string
	Year                  int
	ProjectNameUpper      string
	ProjectIncludeFiles   string
	LibcollectionsInclude string
	LibcollectionsLinker  string
	ProjectNameSnaked     string
}

func CSourceHeader() (*template.Template, error) {
	tmpTpl := template.New("CHeader")
	tpl, err := tmpTpl.Parse(headerContent)

	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func BashSourceHeader() (*template.Template, error) {
	tmpTpl := template.New("BashHeader")
	tpl, err := tmpTpl.Parse(bashHeaderContent)

	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func camelCase(src string) string {
	var camelCase string
	isToUpper := false

	for _, runeValue := range src {
		if isToUpper {
			camelCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			if runeValue == '-' {
				isToUpper = true
			} else {
				camelCase += string(runeValue)
			}
		}
	}

	return camelCase
}

func GetContentData(options base.FileOptions) ContentData {
	now := time.Now()

	return ContentData{
		ProjectName:       options.ProjectName,
		ProjectNameUpper:  strings.ToUpper(options.ProjectName),
		Author:            options.AuthorName,
		Year:              now.Year(),
		Date:              now.Format(time.ANSIC),
		ProjectNameSnaked: camelCase(options.ProjectName),
	}
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

func errorContent(fileOptions ContentType, options base.FileOptions) string {
	if options.LibcollectionsFeatures {
		if fileOptions&Source != 0 {
			return `
static const char *__description[] = {
    cl_tr_noop("Ok"),
};

static const char *__unknown_error = cl_tr_noop("Unknown error");

struct error_storage {
    int error;
};

cl_error_storage_declare(__storage__, sizeof(struct error_storage))
#define __cerrno        (cl_errno_storage(&__storage__))

void errno_clear(void)
{
    struct error_storage *e = __cerrno;

    e->error = {{.ProjectNameUpper}}_NO_ERROR;
}

void errno_set(enum {{.ProjectName}}_error_code code)
{
    struct error_storage *e = __cerrno;

    e->error = code;
}

__PUB_API__ enum {{.ProjectName}}_error_code {{.ProjectName}}_get_last_error(void)
{
    struct error_storage *e = __cerrno;

    return e->error;
}

__PUB_API__ const char *{{.ProjectName}}_strerror(enum {{.ProjectName}}_error_code code)
{
    if (code >= {{.ProjectNameUpper}}_MAX_ERROR_CODE)
        return __unknown_error;

    return __description[code];
}
`
		} else {
			if fileOptions&InternalHeader != 0 {
				return `
enum {{.ProjectName}}_error_code {
    {{.ProjectNameUpper}}_NO_ERROR,

    {{.ProjectNameUpper}}_MAX_ERROR_CODE
};

void errno_clear(void);
void errno_set(enum {{.ProjectName}}_error_code code);
`
			} else {
				return `
enum {{.ProjectName}}_error_code {{.ProjectName}}_get_last_error(void);
const char *{{.ProjectName}}_strerror(enum {{.ProjectName}}_error_code code);
`
			}
		}
	}

	return ""
}

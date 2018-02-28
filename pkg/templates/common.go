package templates

import (
	"strings"
	"text/template"
	"time"

	"source-template/pkg/base"
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
	ProjectName      string
	Author           string
	Date             string
	Year             int
	ProjectNameUpper string
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

func GetContentData(options base.FileOptions) ContentData {
	now := time.Now()

	return ContentData{
		ProjectName:      options.ProjectName,
		ProjectNameUpper: strings.ToUpper(options.ProjectName),
		Author:           options.AuthorName,
		Year:             now.Year(),
		Date:             now.Format(time.ANSIC),
	}
}

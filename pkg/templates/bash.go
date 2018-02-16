package templates

import (
	"fmt"
	"os"

	"source-template/pkg/base"
)

type BashFile struct {
	base.FileOptions
}

func (s BashFile) Header(file *os.File) {
	// if we're creating a project, probably will have an include directive here
}

func (s BashFile) HeaderComment(file *os.File) {
	file.WriteString(`
#
# Description:
#
# Author:
# Created at:
# Project:
#
# Copyright (C) 2017 Author Name All rights reserved.
#
`)
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
	}
}

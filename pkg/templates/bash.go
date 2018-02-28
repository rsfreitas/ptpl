package templates

import (
	"fmt"
	"os"

	"source-template/pkg/base"
)

type BashFile struct {
	base.FileOptions
	ContentData
}

func (s BashFile) Header(file *os.File) {
	// if we're creating a project, probably will have an include directive here
}

func (s BashFile) HeaderComment(file *os.File) {
	tpl, err := BashSourceHeader()

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
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
		ContentData: GetContentData(options),
	}
}

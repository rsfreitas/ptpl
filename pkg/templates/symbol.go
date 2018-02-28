package templates

import (
	"os"
	"text/template"

	"source-template/pkg/base"
)

type SymbolFile struct {
	base.FileOptions
	ContentData
}

const content = `LIB{{.ProjectNameUpper}}_0.1 {
	global:
		*;
	local:
		*;
};
`

func (s SymbolFile) Header(file *os.File) {
}

func (s SymbolFile) HeaderComment(file *os.File) {
}

func (s SymbolFile) Footer(file *os.File) {
}

func (s SymbolFile) Content(file *os.File) {
	tmpTpl := template.New("symbol")
	tpl, err := tmpTpl.Parse(content)

	if err != nil {
		return
	}

	tpl.Execute(file, s.ContentData)
}

func NewSymbol(options base.FileOptions) base.FileTemplate {
	return &SymbolFile{
		FileOptions: options,
		ContentData: GetContentData(options),
	}
}

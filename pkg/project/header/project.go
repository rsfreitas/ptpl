package header

import (
	"fmt"

	"source-template/pkg/base"
	"source-template/pkg/templates"
)

type SingleHeader struct {
	file base.FileInfo
	base.ProjectOptions
}

func (s SingleHeader) String() string {
	return fmt.Sprintf("Project type: single header\nFilename: %s\n", s.ProjectName)
}

func (s SingleHeader) Build() error {
	return s.file.Build()
}

func New(options base.ProjectOptions) (base.Project, error) {
	fileOptions := base.FileOptions{
		Name:           base.AddExtension(options.ProjectName, ".h"),
		HeaderComment:  true,
		ProjectOptions: options,
	}

	return &SingleHeader{
		file: base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewHeader(fileOptions, ""),
		},
		ProjectOptions: options,
	}, nil
}

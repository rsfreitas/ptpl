package xante

import (
	"source-template/pkg/base"
)

type XantePlugin struct {
	sources  []base.FileInfo
	headers  []base.FileInfo
	rootPath string
	base.ProjectOptions
}

func (l XantePlugin) String() string {
	return ""
}

func (l XantePlugin) Build() error {
	// create root path and subdirs
	// create sources
	// create headers
	// create Makefile (future CMakeLists.txt)
	// create application script

	return nil
}

func New(options base.ProjectOptions) (base.Project, error) {
	return &XantePlugin{}, nil
}

package library

import (
	"fmt"
	"os"

	"source-template/pkg/base"
	"source-template/pkg/project/common"
	"source-template/pkg/templates"
)

type Library struct {
	sources  []base.FileInfo
	headers  []base.FileInfo
	debian   []base.FileInfo
	rootPath string
	base.ProjectOptions
}

func (l Library) String() string {
	return fmt.Sprintf("Library project")
}

func createLibraryDirtree(path string, options base.ProjectOptions) error {
	var subdirs []string
	var prefix string

	if options.PackageProject {
		prefix = options.ProjectName
		subdirs = append(subdirs, "pkg_install/misc")
		subdirs = append(subdirs, "pkg_install/debian")
	}

	subdirs = append(subdirs, prefix+"/src")
	subdirs = append(subdirs, prefix+"/include")

	for _, dir := range subdirs {
		err := os.MkdirAll(path+"/"+dir, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}

func (l Library) Build() error {
	// create root path and subdirs
	if err := createLibraryDirtree(l.rootPath, l.ProjectOptions); err != nil {
		return err
	}

	// create sources
	for _, f := range l.sources {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create headers
	for _, f := range l.headers {
		if err := f.Build(); err != nil {
			return err
		}
	}

	// create Makefile (future CMakeLists.txt)

	return nil
}

func createSources(options base.ProjectOptions, rootPath string, prefix string) []base.FileInfo {
	var files []base.FileInfo

	sources := []string{
		"utils",
		"log",
		"error",
	}

	for _, s := range sources {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(rootPath+"/"+prefix+"/src/"+s, ".c"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewSource(fileOptions),
		})
	}

	return files
}

func createHeaders(options base.ProjectOptions, rootPath string, prefix string) []base.FileInfo {
	var files []base.FileInfo
	var headers []string

	headers = append(headers, "lib"+options.ProjectName)
	headers = append(headers, "lib"+options.ProjectName+"_internal.h")

	for _, h := range headers {
		fileOptions := base.FileOptions{
			ProjectOptions: options,
			HeaderComment:  true,
			Name:           base.AddExtension(rootPath+"/"+prefix+"/include/"+h, ".h"),
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewHeader(fileOptions, ""),
		})
	}

	return files
}

func New(options base.ProjectOptions) (base.Project, error) {
	var rootPath string
	var prefix string
	cwd, err := os.Getwd()

	if err != nil {
		return &Library{}, err
	}

	if options.PackageProject {
		prefix = options.ProjectName
		rootPath = cwd + "/package-lib" + options.ProjectName
	} else {
		rootPath = cwd + "/lib" + options.ProjectName
	}

	return &Library{
		rootPath:       rootPath,
		sources:        createSources(options, rootPath, prefix),
		headers:        createHeaders(options, rootPath, prefix),
		debian:         common.CreateDebianScripts(options, rootPath, prefix),
		ProjectOptions: options,
	}, nil
}

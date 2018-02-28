package common

import (
	"source-template/pkg/base"
	"source-template/pkg/templates"
)

func CreateDebianScripts(options base.ProjectOptions, rootPath string) []base.FileInfo {
	var files []base.FileInfo
	scripts := []string{
		"preinst",
		"prerm",
		"postinst",
		"postrm",
	}

	// If we're not a package
	if !options.PackageProject {
		return files
	}

	for _, s := range scripts {
		fileOptions := base.FileOptions{
			Executable:     true,
			HeaderComment:  true,
			ProjectOptions: options,
			Name:           rootPath + "/pkg_install/debian/" + s,
		}

		files = append(files, base.FileInfo{
			FileOptions:  fileOptions,
			FileTemplate: templates.NewBash(fileOptions),
		})
	}

	return files
}

func CreateMakefile(options base.ProjectOptions, rootPath string, prefix string) base.FileInfo {
	fileOptions := base.FileOptions{
		Executable:     false,
		HeaderComment:  false,
		ProjectOptions: options,
		Name:           rootPath + "/" + prefix + "/CMakeLists.txt",
	}

	return base.FileInfo{
		FileOptions:  fileOptions,
		FileTemplate: templates.NewMakefile(fileOptions),
	}
}

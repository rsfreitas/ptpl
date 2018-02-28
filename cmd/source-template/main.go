// source-template is an application to create source project templates to
// improve software development.
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
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"source-template/pkg/base"
	"source-template/pkg/project"
)

const AppName string = "source-tpl"
const Version string = "0.2.0"

type CLIOptions struct {
	version bool
	quiet   bool
	base.ProjectOptions
}

func validateProjectType(projectType int) error {
	switch {
	case projectType >= base.SingleSourceProject && projectType <= base.XantePluginProject:
		return nil
	}

	return errors.New("Unsupported chosen project")
}

func validateProjectLanguage(language int, projectType int) error {
	if language != base.CLanguage && projectType != base.XantePluginProject {
		return errors.New("Programming language unsupported for this kind of project")
	}

	switch {
	case language >= base.CLanguage && language <= base.RustLanguage:
		return nil
	}

	return errors.New("Unsupported programming language")
}

// validateOptions does the command line options validations
func validateOptions(options CLIOptions) error {
	err := validateProjectType(options.ProjectType)

	if err != nil {
		return err
	}

	err = validateProjectLanguage(options.Language, options.ProjectType)

	if err != nil {
		return err
	}

	return nil
}

// TODO: Add description of options
// getCLIOptions configures the application supported command line options.
func getCLIOptions() CLIOptions {
	var options CLIOptions
	var projectType, language string

	flag.BoolVar(&options.version, "v", false,
		"Shows the current application version.")

	flag.BoolVar(&options.PackageProject, "package", false,
		"Enables template creation as a project.")

	flag.StringVar(&options.ProjectName, "name", "",
		"Assigns the project's name.")

	language, err := base.LanguageKey(base.CLanguage)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	flag.StringVar(&language, "language", language,
		"Chooses the programming language to the created template.")

	flag.StringVar(&options.AuthorName, "author", "",
		"Assigns the project author's name.")

	defaultProject, err := base.ProjectKey(base.SingleSourceProject)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	flag.StringVar(&projectType, "type", defaultProject,
		"Chooses the template project type.")

	flag.BoolVar(&options.quiet, "quiet", false,
		"Disables project creation messages.")

	flag.Parse()

	if options.version {
		fmt.Printf("%s - version %s\n", AppName, Version)
		os.Exit(0)
	}

	options.ProjectType, err = base.ProjectLookup(projectType)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	options.Language, err = base.LanguageLookup(language)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err := validateOptions(options); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	return options
}

func main() {
	options := getCLIOptions()
	p, err := project.Assemble(options.ProjectOptions)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if !options.quiet {
		fmt.Println(p)
	}

	if err := p.Build(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

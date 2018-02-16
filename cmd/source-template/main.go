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
const Version string = "0.1.0"

type CLIOptions struct {
	version bool
	quiet   bool
	base.ProjectOptions
}

//buildCLIOptions configures the application supported command line options.
func buildCLIOptions(options *CLIOptions) {
	flag.BoolVar(&options.version, "v", false,
		"Shows the current application version.")

	flag.BoolVar(&options.PackageProject, "package", false,
		"Enables template creation as a project.")

	flag.StringVar(&options.ProjectName, "name", "",
		"Assigns the project's projectName.")

	flag.IntVar(&options.Language, "language", base.CLanguage,
		"Chooses the programming language to the created template.")

	flag.StringVar(&options.AuthorName, "author", "",
		"ASsigns the project author's projectName.")

	flag.IntVar(&options.ProjectType, "type", base.SingleSourceProject,
		"Chooses the template project type.")

	flag.BoolVar(&options.quiet, "quiet", false,
		"Disables project creation messages.")
}

func validateProjectType(projectType int) error {
	switch {
	case projectType >= base.SingleSourceProject && projectType <= base.XantePluginProject:
		return nil
	}

	return errors.New("Unsupported chosen project")
}

func validateProjectLanguage(language int) error {
	switch {
	case language >= base.CLanguage && language <= base.RustLanguage:
		return nil
	}

	return errors.New("Unsupported programming language")
}

//validateOptions does the command line options validations
func validateOptions(options CLIOptions) error {
	if err := validateProjectType(options.ProjectType); err != nil {
		return err
	}

	if err := validateProjectLanguage(options.Language); err != nil {
		return err
	}

	return nil
}

func main() {
	var options CLIOptions

	buildCLIOptions(&options)
	flag.Parse()

	if options.version {
		fmt.Printf("%s - version %s\n", AppName, Version)
		os.Exit(0)
	}

	if err := validateOptions(options); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

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

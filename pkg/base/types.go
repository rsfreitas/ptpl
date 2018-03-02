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
package base

import (
	"errors"
)

const (
	SingleSourceProject = 1 + iota
	SingleHeaderProject
	ApplicationProject
	LibraryProject
	XantePluginProject
)

const (
	CLanguage = 1 + iota
	JavaLanguage
	PythonLanguage
	GoLanguage
	RustLanguage
)

var supportedProjects = map[string]int{
	"header":       SingleHeaderProject,
	"source":       SingleSourceProject,
	"application":  ApplicationProject,
	"library":      LibraryProject,
	"xante-plugin": XantePluginProject,
}

var supportedLanguages = map[string]int{
	"C":      CLanguage,
	"java":   JavaLanguage,
	"python": PythonLanguage,
	"go":     GoLanguage,
	"rust":   RustLanguage,
}

func ProjectLookup(project string) (int, error) {
	code := supportedProjects[project]

	if code == 0 {
		return -1, errors.New("Unknown project")
	}

	return code, nil
}

func ProjectKey(project int) (string, error) {
	for k, v := range supportedProjects {
		if v == project {
			return k, nil
		}
	}

	return "", errors.New("Unknown project")
}

func LanguageLookup(language string) (int, error) {
	code := supportedLanguages[language]

	if code == 0 {
		return -1, errors.New("Unknown language")
	}

	return code, nil
}

func LanguageKey(language int) (string, error) {
	for k, v := range supportedLanguages {
		if v == language {
			return k, nil
		}
	}

	return "", errors.New("Unknown language")
}

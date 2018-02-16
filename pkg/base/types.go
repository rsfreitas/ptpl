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

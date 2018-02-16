package base

import (
	"path/filepath"
	"strings"
)

func AddExtension(filename string, extension string) string {
	ext := filepath.Ext(filename)

	if ext != "" {
		return filename
	}

	fileExtension := extension

	if !strings.Contains(fileExtension, ".") {
		fileExtension = fileExtension[1:]
	}

	return filename + fileExtension
}

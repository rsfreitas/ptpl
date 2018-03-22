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
	"os"
	"os/exec"
)

type FileTemplate interface {
	Header(file *os.File)
	HeaderComment(file *os.File)
	Footer(file *os.File)
	Content(file *os.File)
}

type FileOptions struct {
	Name          string
	HeaderComment bool
	Executable    bool
	LibraryHeader bool
	ProjectOptions
}

type FileInfo struct {
	FileOptions
	FileTemplate
}

func (f FileInfo) Build(path string) error {
	filename := path + "/" + f.Name
	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	if f.FileOptions.HeaderComment {
		f.FileTemplate.HeaderComment(file)
	}

	f.Header(file)
	f.Content(file)
	f.Footer(file)

	if f.Executable {
		cmd := exec.Command("chmod", "+x", filename)

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func New() *FileInfo {
	return &FileInfo{}
}

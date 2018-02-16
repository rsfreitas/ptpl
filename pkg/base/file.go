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

func (f FileInfo) Build() error {
	file, err := os.Create(f.Name)

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
		cmd := exec.Command("chmod", "+x", f.Name)

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func New() *FileInfo {
	return &FileInfo{}
}

package parser

import (
	"fmt"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"path/filepath"
)

func GetModule(path string) (modPath string, file *modfile.File, err error) {
	for {
		goModPath := filepath.Join(path, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			modPath = goModPath
			break
		}
		path = filepath.Dir(path)
		if path == "." {
			return "", nil, fmt.Errorf("go.mod file not exist")
		}
	}
	f, err := os.Open(modPath)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	filename := filepath.Base(modPath)
	file, err = modfile.ParseLax(filename, data, func(path, version string) (string, error) {
		return version, nil
	})
	if err != nil {
		return "", nil, err
	}
	return
}

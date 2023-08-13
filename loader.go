package genpig

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func jsonLoadFromFilePath(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if json.NewDecoder(f).Decode(v) == nil {
		return nil
	}
	return fmt.Errorf("not found")
}

func getConfigPaths(dirs, filenames []string) (configFiles []string) {
	for _, searchPath := range dirs {
		searchPath = filepath.Join(os.ExpandEnv(searchPath))
		filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() && path != searchPath {
				return filepath.SkipDir
			}
			for _, filename := range filenames {
				if strings.TrimSuffix(d.Name(), filepath.Ext(d.Name())) == filename {
					if strings.ToLower(filepath.Ext(d.Name())) == ".json" {
						configFiles = append(configFiles, searchPath)
					}
				}
			}
			return nil
		})
	}
	return
}

func LoadJsonConfig(dirs, filenames []string, v any) (loaded bool) {
	filePaths := getConfigPaths(dirs, filenames)
	if len(filePaths) != 0 {
		for _, path := range filePaths {
			if jsonLoadFromFilePath(path, v) == nil {
				loaded = true
				break
			}
		}
	}
	return
}

func Merge[T any](values ...T) (merged T) {
	for _, v := range values {
		if !reflect.ValueOf(v).IsZero() {
			merged = v
		}
	}
	return
}

func EnvFloat(key string) (value float64) {
	envValue := os.Getenv(key)
	if envValue == "" {
		return
	}
	value, _ = strconv.ParseFloat(envValue, 64)
	return
}

func EnvInt(key string) (value int64) {
	envValue := os.Getenv(key)
	if envValue == "" {
		return
	}
	value, _ = strconv.ParseInt(envValue, 10, 64)
	return
}

func EnvBool(key string) (value bool) {
	envValue := os.Getenv(key)
	if envValue == "" {
		return
	}
	value, _ = strconv.ParseBool(envValue)
	return
}

package site

import (
	"io/fs"
	"path/filepath"
)

type Site struct {
	Volumes map[string]string `toml:"volumes"`
}

func (s *Site) ScanAll(root string) error {

	var pagePaths []string

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".pagine" {
			pagePaths = append(pagePaths, path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

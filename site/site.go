// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package site

import (
	. "github.com/webpagine/pagine/page"
	"github.com/webpagine/pagine/util"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Site struct {
	Root string
}

func (s *Site) ScanAll(path string) (pagePaths []string, err error) {
	err = filepath.Walk(filepath.Join(s.Root, path), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".pagine" {
			pagePaths = append(pagePaths, path)
		}

		return nil
	})
	if err != nil {
		return
	}

	return
}

func (s *Site) GenerateAll(dst string) error {
	pagePaths, err := s.ScanAll(".")
	if err != nil {
		return err
	}

	var pages []Page

	for _, pagePath := range pagePaths {
		var page Page

		err = util.UnmarshalTOMLFile(pagePath, &page)
		if err != nil {
			return err
		}

		page.Path = pagePath

		pages = append(pages, page)
	}

	for _, page := range pages {
		outputPath, _ := strings.CutSuffix(page.Path, ".pagine")

		absolutePath := filepath.Join(dst, outputPath)

		err = os.MkdirAll(filepath.Dir(absolutePath), os.ModePerm)
		if err != nil {
			return err
		}

		outputFile, err := os.Create(absolutePath)
		if err != nil {
			return err
		}

		err = page.Generate(outputFile)
		if err != nil {
			return err
		}
	}

	return nil
}

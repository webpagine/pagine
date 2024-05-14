// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package site

import (
	. "github.com/webpagine/go-pagine/page"
	. "github.com/webpagine/go-pagine/path"
	"github.com/webpagine/go-pagine/util"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Site struct {
	Root Path
}

func (s *Site) GenerateAll(dest string) error {

	var (
		pagePathList   []string
		staticPathList []string
	)

	pathList, err := s.Root.IterateFilesAsRelative()
	if err != nil {
		return err
	}

	for _, path := range pathList {
		if strings.HasSuffix(path, ".pagine") {
			pagePathList = append(pagePathList, path)
		} else {
			staticPathList = append(staticPathList, path)
		}
	}

	errs := make(chan error, len(pathList))
	var wg sync.WaitGroup

	var pages []Page

	for _, pagePath := range pagePathList {
		var page Page

		err = util.UnmarshalTOMLFile(s.Root.PathOf(pagePath), &page)
		if err != nil {
			return err
		}

		page.Path = pagePath

		pages = append(pages, page)
	}

	for _, staticPath := range staticPathList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = util.CopyFile(s.Root.PathOf(staticPath), filepath.Join(dest, staticPath))
			if err != nil {
				errs <- err
			}
		}()
	}

	// Generate pages.

	for _, page := range pages {

		// Destination file name.
		outputRelativePath, _ := strings.CutSuffix(page.Path, ".pagine")
		outputAbsolutePath := s.Root.PathOf(filepath.Join(dest, outputRelativePath))

		// Create dest file.

		_, err := os.Stat(filepath.Dir(outputAbsolutePath))
		switch {
		case os.IsNotExist(err):
			err = os.MkdirAll(filepath.Dir(outputAbsolutePath), os.ModePerm)
			if err != nil {
				return err
			}
		case err == nil:
			return err
		}

		outputFile, err := os.Create(outputAbsolutePath)
		if err != nil {
			return err
		}

		// Generate in parallel.
		wg.Add(1)
		go func() {
			defer outputFile.Close()
			defer wg.Done()
			err := page.Generate(outputFile)
			if err != nil {
				errs <- err
			}
		}()
	}

	wg.Wait()

	close(errs)

	var errSet util.ErrorSet

	for {
		err, ok := <-errs
		if !ok {
			break
		}
		errSet.Errors = append(errSet.Errors, err)
	}

	if errSet.Errors != nil {
		return &errSet
	}

	return nil
}

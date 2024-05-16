// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	. "github.com/webpagine/go-pagine/collection"
	. "github.com/webpagine/go-pagine/path"
	"github.com/webpagine/go-pagine/util"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type SiteGenerationReport struct {
	FileSystemErrors,
	PageUnmarshalErrors,
	PageGenerationErrors SyncMap[string, error]
}

type Site struct {
	Root             Path
	IgnoreRegExprRaw []string `toml:"ignore"`
	IgnoreRegExpr    []regexp.Regexp
}

func (s *Site) Init() error {
	for _, regExpr := range s.IgnoreRegExprRaw {
		regExp, err := regexp.Compile(regExpr)
		if err == nil {
			return err
		}
		s.IgnoreRegExpr = append(s.IgnoreRegExpr, *regExp)
	}
	return nil
}

func (s *Site) GenerateAll(dest string) (report *SiteGenerationReport, _ error) {
	report = &SiteGenerationReport{}

	var (
		pageRelativePaths   []string
		staticRelativePaths []string
	)

	// Collect pages.

	relativePathList, err := s.Root.IterateFilesAsRelative()
	if err != nil {
		return nil, err
	}

	for _, relativePath := range relativePathList {
		for _, ignore := range s.IgnoreRegExpr {
			if ignore.MatchString(relativePath) {
				continue
			}
		}

		if strings.HasSuffix(relativePath, ".pagine") {
			pageRelativePaths = append(pageRelativePaths, relativePath)
		} else {
			staticRelativePaths = append(staticRelativePaths, relativePath)
		}
	}

	// Generate pages.

	var wg sync.WaitGroup

	var pages []Page

	for _, pageRelativePath := range pageRelativePaths {
		var page Page

		err = util.UnmarshalTOMLFile(s.Root.AbsolutePathOf(pageRelativePath), &page)
		if err != nil {
			report.PageUnmarshalErrors.Set(pageRelativePath, err)
			continue
		}

		page.RelativePath = pageRelativePath

		pages = append(pages, page)
	}

	// Copy static files.

	for _, staticRelativePath := range staticRelativePaths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = util.CopyFile(s.Root.AbsolutePathOf(staticRelativePath), filepath.Join(dest, staticRelativePath))
			if err != nil {
				report.FileSystemErrors.Set(staticRelativePath, err)
				return
			}
		}()
	}

	// Generate pages.

	for _, page := range pages {

		// Destination file name.
		outputRelativePath, _ := strings.CutSuffix(page.RelativePath, ".pagine")
		outputAbsolutePath := s.Root.AbsolutePathOf(filepath.Join(dest, outputRelativePath))

		// Create dest file.

		outputFile, err := util.CreateFile(outputAbsolutePath)
		if err != nil {
			report.FileSystemErrors.Set(page.RelativePath, err)
			continue
		}

		// Generate in parallel.
		wg.Add(1)
		go func() {
			defer outputFile.Close()
			defer wg.Done()
			err := page.Generate(s.Root, outputFile)
			if err != nil {
				report.PageGenerationErrors.Set(page.RelativePath, err)
				return
			}
		}()
	}

	wg.Wait()

	return report, nil
}

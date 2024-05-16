// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"flag"
	"fmt"
	. "github.com/webpagine/go-pagine/path"
	. "github.com/webpagine/go-pagine/structure"
	"github.com/webpagine/go-pagine/util"
	"os"
	"path/filepath"
)

var (
	wd, _      = os.Getwd()
	doGenerate = flag.Bool("gen", false, "GenerateAll site.")
	doServe    = flag.Bool("serve", false, "Serve as HTTP.")
	siteRoot   = flag.String("root", wd, "Site root.")
	publicDir  = flag.String("public", "../"+filepath.Base(wd)+".public", "Location of generated site.")
)

func main() {

	flag.Parse()

	err := _main()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func _main() error {

	if *siteRoot == "" {
		return errors.New("site root is required")
	}

	var site = Site{
		Root: NewPath(*siteRoot),
	}

	err := site.Init()
	if err != nil {
		return err
	}

	err = util.UnmarshalTOMLFile(site.Root.AbsolutePathOf("/pagine.toml"), &site)
	if err != nil {
		return err
	}

	if *doServe {
		// TODO
	}

	if *doGenerate {
		report, err := site.GenerateAll(*publicDir)
		if err != nil {
			return err
		}

		if report.FileSystemErrors.RawMap != nil {
			fmt.Println("File system errors:")
			for _, err := range report.FileSystemErrors.RawMap {
				fmt.Print("\t", err, "\n")
			}
		}

		if report.PageUnmarshalErrors.RawMap != nil {
			fmt.Println("Page unmarshal errors:")
			for _, err := range report.PageUnmarshalErrors.RawMap {
				fmt.Print("\t", err, "\n")
			}
		}

		if report.PageGenerationErrors.RawMap != nil {
			fmt.Println("Page generation errors:")
			for _, err := range report.PageGenerationErrors.RawMap {
				fmt.Print("\t", err, "\n")
			}
		}
	}

	return nil
}

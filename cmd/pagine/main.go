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
	wd, _ = os.Getwd()

	optGenerate = flag.Bool("gen", false, "GenerateAll site.")
	optServe    = flag.Bool("serve", false, "Serve as HTTP.")

	siteRoot  = flag.String("root", wd, "Site root.")
	publicDir = flag.String("public", "/tmp/"+filepath.Base(wd)+".public", "Location of generated site.")

	addr = flag.String("listen", ":8080", "Listen address.")
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

	var site = &Site{
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

	if *optGenerate {
		err := doGenerate(site)
		if err != nil {
			return err
		}
	}

	if *optServe {
		err := doServe(site)
		if err != nil {
			return err
		}
	}

	return nil
}

func doGenerate(site *Site) error {
	report, err := site.GenerateAll(*publicDir)
	if err != nil {
		return err
	}

	PrintReport(report)

	return nil
}

func doServe(site *Site) error {
	err := doGenerate(site)
	if err != nil {
		return err
	}

	err = Serve(*addr, site, *publicDir)
	if err != nil {
		return err
	}

	return nil
}

// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"github.com/webpagine/pagine/v2/vfs"
	"os"
	"path/filepath"
)

var (
	wd, _ = os.Getwd()

	optRootDir   = flag.String("root", wd, "Site root.")
	optPublicDir = flag.String("public", "/tmp/"+filepath.Base(wd)+".public", "Location of public directory.")

	optAddr = flag.String("serve", "", "Listen and Serve as HTTP.")

	jobBuilderRoot = os.Getenv("PAGINE_JOB_BUILDER_ROOT")
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
	root := vfs.OsDirFS(*optRootDir)
	dest := vfs.OsDirFS(*optPublicDir)

	if *optAddr != "" {
		err := Serve(root, dest)
		if err != nil {
			return err
		}
	} else {
		err := GenerateAll(root, dest, jobBuilderRoot, false)
		if err != nil {
			return err
		}
	}

	return nil
}

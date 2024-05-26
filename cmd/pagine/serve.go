// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	. "github.com/webpagine/go-pagine/structure"
	"log"
	"net/http"
	"os"
)

func Serve(addr string, site *Site, publicDir string) error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	fileRelativePaths, err := site.Root.IterateFilesAsRelative()
	if err != nil {
		return err
	}

	for _, relativePath := range fileRelativePaths {
		err = watcher.Add(site.Root.AbsolutePathOf(relativePath))
		if err != nil {
			return err
		}
	}

	done := make(chan struct{}, 1)
	defer close(done)

	// Handle file changes.
	go func() {
		gen := func() {
			err := func() error {
				report, err := site.GenerateAll(publicDir)
				if err != nil {
					return err
				}

				PrintReport(report)

				return nil
			}()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Println("Executed full generation.")
		}
		gen()

		for {
			select {
			case <-watcher.Events:
				// TODO Incremental generation.

				// Full generation.

				log.Println("Detected file change.")

				gen()
			case <-done:
				return
			}
		}
	}()

	err = http.ListenAndServe(addr, http.FileServer(http.Dir(publicDir)))
	if err != nil {
		return err
	}

	return nil
}

// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package environment

import (
	"github.com/fsnotify/fsnotify"
	. "github.com/webpagine/pagine/site"
	"net/http"
)

func Serve(site *Site, publicDir string) error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(site.Root)
	if err != nil {
		return err
	}

	done := make(chan struct{}, 1)

	// Handle file changes.
	go func() {
		select {
		case <-watcher.Events:
			err = site.GenerateAll(publicDir)
			if err != nil {
				return
			}
		case <-done:
			break
		}
	}()

	err = http.ListenAndServe("/", http.FileServer(http.Dir(publicDir)))
	if err != nil {
		return err
	}

	return nil
}

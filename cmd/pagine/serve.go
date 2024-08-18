// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/webpagine/pagine/v2/vfs"
	"io/fs"
	"net/http"
	"sync/atomic"
	"time"
)

func Serve(root, dest *vfs.DirFS) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	notify := func() error {
		return fs.WalkDir(root, "/", func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				sub, err := root.Chroot(path)
				if err != nil {
					return err
				}

				err = w.Add(sub.Path)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	err = notify()
	if err != nil {
		return err
	}

	genEvent := make(chan struct{})

	var updated atomic.Bool

	go func() {
		for {
			if !updated.Load() {
				updated.Store(true)
				err := GenerateAll(root, dest, true)
				if err != nil {
					fmt.Println(err)
				}
				close(genEvent)
				genEvent = make(chan struct{})
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case e := <-w.Events:
				fmt.Println(e)
				updated.Store(false)
			}

			err := notify()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(dest.Path)))
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var upgrader websocket.Upgrader
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		go func() {
			<-genEvent
			_ = conn.Close()
		}()
	})

	return http.ListenAndServe(*optAddr, mux)
}

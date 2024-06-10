// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/webpagine/pagine/structure"
	"github.com/webpagine/pagine/vfs"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	wd, _ = os.Getwd()

	optRootDir   = flag.String("root", wd, "Site root.")
	optPublicDir = flag.String("public", "/tmp/"+filepath.Base(wd)+".public", "Location of public directory.")

	optAddr = flag.String("serve", "", "Listen and serve as HTTP.")
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

	err := generateAll(root, dest)
	if err != nil {
		return err
	}

	if *optAddr != "" {
		err = serve(root, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateAll(root, dest vfs.DirFS) error {
	env, err := structure.LoadEnv(root)
	if err != nil {
		fmt.Println("Error occurred while loading environment from env.toml:")
		return err
	}

	err = fs.WalkDir(root, "/", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		for _, glob := range env.IgnoreGlobs {
			if glob.MatchString(path) {
				return nil
			}
		}

		src, err := root.Open(path)
		if err != nil {
			return err
		}

		dst, err := dest.CreateFile(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error occurred while copying files:")
		return err
	}

	level, err := structure.ExecuteLevels(env, root, dest, structure.MetadataSet{})
	if err != nil {
		fmt.Println("Error occurred while executing units:")
		return err
	}

	walkLevels(&level)

	fmt.Println("Generation complete.")

	return nil
}

func walkLevels(level *structure.Level) {
	for _, report := range level.Reports {
		fmt.Println("[", report.Level.Root.Path, "]", report.Err)
	}

	f := sync.OnceFunc(func() {
		println(level.Root.Path)
	})

	for _, u := range level.Units {
		switch {
		case u.Report.Error != nil:
			f()
			fmt.Println("\t[", u.Output, "]", u.Report.Error)
		case u.Report.TemplateErrors != nil:
			f()
			fmt.Println("\t[", u.Output, "] Template errors")
			for _, e := range u.Report.TemplateErrors {
				fmt.Println("\t\t", e)
			}
		default:
		}
	}

	for _, level := range level.Levels {
		walkLevels(&level)
	}
}

func serve(root, dest vfs.DirFS) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = w.Add(root.Path)
	if err != nil {
		return err
	}
	defer w.Close()

	go func() {
		select {
		case <-w.Events:
			err := generateAll(root, dest)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	return http.ListenAndServe(*optAddr, http.FileServerFS(dest))
}

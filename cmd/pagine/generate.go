// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/jellyterra/go-fscopy"
	"github.com/webpagine/pagine/v2/config"
	"github.com/webpagine/pagine/v2/structure"
	"github.com/webpagine/pagine/v2/vfs"
	"io/fs"
	"os"
	"time"
)

func GenerateAll(root, dest *vfs.DirFS, jobBuilderRoot string, isServing bool) error {

	env, err := config.LoadEnv(root)
	if err != nil {
		fmt.Println("Error occurred while loading environment from env.toml:")
		return err
	}

	env.DeployTime = time.Now().UTC().Format("20060102150405")
	env.IsServing = isServing

	err = os.RemoveAll(dest.Path)
	switch {
	case err == nil:
	case os.IsNotExist(err):
	default:
		return err
	}

	err = fscopy.CopyAllWithExceptionGlobs(root.Path, dest.Path, env.IgnoreGlobs...)
	if err != nil {
		fmt.Println("Error occurred while copying files:")
		return err
	}

	level, err := config.CollectLevels(env, env.Root, structure.MetadataSet{})
	if err != nil {
		return err
	}

	ExecuteLevels(env, dest, level)

	err = fs.WalkDir(dest, "/", func(path string, d fs.DirEntry, err error) error {
		return nil
	})
	if err != nil {
		return err
	}

	if jobBuilderRoot == "" {
		err = CollectAndRunWorkflows(root, dest, nil)
	} else {
		err = CollectAndRunWorkflows(root, dest, os.DirFS(jobBuilderRoot).(fs.StatFS))
	}
	if err != nil {
		return err
	}

	fmt.Println("Generation complete.")

	return nil
}

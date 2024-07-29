// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package vfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func OsDirFS(path string) *DirFS {
	return &DirFS{Path: path}
}

type DirFS struct {
	Path string
}

func (dir *DirFS) Chroot(path string) (*DirFS, error) {
	err := dir.Validate("Chroot", path)
	if err != nil {
		return &DirFS{}, err
	}

	sub := &DirFS{Path: filepath.Join(dir.Path, path)}

	stat, err := os.Stat(sub.Path)
	if err != nil {
		return &DirFS{}, err
	}

	if !stat.IsDir() {
		return &DirFS{}, &PathError{
			Op:   "Chroot",
			Path: sub.Path,
			Err:  fmt.Errorf("%s is not a directory", path),
		}
	}

	return sub, nil
}

func (dir *DirFS) Validate(op, name string) error {
	if strings.Contains(name, "..") {
		return &PathError{
			Op:   op,
			Path: name,
			Err:  &IllegalPathError{},
		}
	}

	return nil
}

func (dir *DirFS) Open(name string) (fs.File, error) {
	err := dir.Validate("Open", name)
	if err != nil {
		return nil, err
	}

	return os.Open(filepath.Join(dir.Path, name))
}

func (dir *DirFS) ReadFile(name string) ([]byte, error) {
	err := dir.Validate("ReadFile", name)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(filepath.Join(dir.Path, name))
}

func (dir *DirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	err := dir.Validate("ReadDir", name)
	if err != nil {
		return nil, err
	}

	return os.ReadDir(filepath.Join(dir.Path, name))
}

func (dir *DirFS) Stat(name string) (fs.FileInfo, error) {
	err := dir.Validate("Stat", name)
	if err != nil {
		return nil, err
	}

	return os.Stat(filepath.Join(dir.Path, name))
}

func (dir *DirFS) CreateFile(name string) (*os.File, error) {
	err := dir.Validate("CreateFile", name)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir.Path, name)

	parentDir := filepath.Dir(path)
	_, err = os.Stat(parentDir)
	switch {
	case err == nil:
	case os.IsNotExist(err):
		err = os.MkdirAll(parentDir, 0755)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	return os.Create(path)
}

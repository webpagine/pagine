// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package path

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func (p *Path) IterateAsRelative() (pathList []string, err error) {
	err = filepath.Walk(p.Path, func(path string, info fs.FileInfo, err error) error {
		path, _ = strings.CutPrefix(path, p.Path)
		pathList = append(pathList, path)
		return nil
	})

	return
}

func (p *Path) IterateFilesAsRelative() (pathList []string, err error) {
	err = filepath.Walk(p.Path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		path, _ = strings.CutPrefix(path, p.Path)
		pathList = append(pathList, path)
		return nil
	})

	return
}

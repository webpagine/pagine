// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package path

import "path/filepath"

type Path struct {
	Path string
}

func NewPath(path string) Path { return Path{Path: path} }

func (p *Path) PathOf(path string) string { return filepath.Join(p.Path, path) }

// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package render

import (
	"github.com/webpagine/pagine/v2/vfs"
	"io"
)

type Renderer func(content []byte) (string, error)

// Renderers
// Register renderers in independent packages by init()
var Renderers = map[string]Renderer{}

func FromFile(r Renderer, file io.Reader) (string, error) {
	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return r(b)
}

func FromPath(r Renderer, root *vfs.DirFS, path string) (string, error) {
	f, err := root.Open(path)
	if err != nil {
		return "", err
	}

	return FromFile(r, f)
}

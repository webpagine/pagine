// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"github.com/webpagine/go-pagine/render"
	"os"
	"path/filepath"
)

type Content struct {
	AbsolutePath string
}

func NewContent(absolutePath string) Content {
	return Content{AbsolutePath: absolutePath}
}

func (c *Content) Generate() ([]byte, error) {
	renderer, found := render.Renderers[filepath.Ext(c.AbsolutePath)]
	if !found {
		return nil, &render.NoRendererFoundError{Path: c.AbsolutePath}
	}

	f, err := os.Open(c.AbsolutePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return renderer(f)
}

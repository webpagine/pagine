// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package render

import (
	"fmt"
	"io"
)

type NoRendererFoundError struct {
	Path string
}

func (e *NoRendererFoundError) Error() string {
	return fmt.Sprint("no renderer found for the content type: ", e.Path)
}

type Renderer func(r io.Reader) ([]byte, error)

// Renderers
// Register renderers in independent packages by init()
var Renderers = map[string]Renderer{}

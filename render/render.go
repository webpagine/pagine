// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package render

import "io"

type Renderer func(r io.Reader, w io.Writer) error

// Renderers
// Register renderers in independent packages by init()
var Renderers = map[string]Renderer{}

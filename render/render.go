package render

import "io"

type Renderer func(r io.Reader, w io.Writer) error

// Renderers
// Register renderers in independent packages by init()
var Renderers = map[string]Renderer{}

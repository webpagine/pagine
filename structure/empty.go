// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

var emptyFuncMap = map[string]any{
	"add": add,
	"sub": sub,
	"mul": mul,
	"div": div,
	"mod": mod,

	"hasPrefix":  hasPrefix,
	"trimPrefix": trimPrefix,

	"getAttr": empty,
	"getEnv":  empty,

	"getMetadata": empty,

	"apply":          __empty,
	"embed":          _empty,
	"render":         _empty,
	"renderMarkdown": _empty,
	"renderAsciidoc": _empty,
}

func empty() any { return "" }

func _empty(_ any) any { return "" }

func __empty(_ any, _ any) any { return "" }

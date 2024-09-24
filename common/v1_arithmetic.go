// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package common

var V1Arithmetic = map[string]any{
	"add": func(aInt, bInt int) int { return aInt + bInt },
	"sub": func(aInt, bInt int) int { return aInt - bInt },
	"mul": func(aInt, bInt int) int { return aInt * bInt },
	"div": func(aInt, bInt int) int { return aInt / bInt },
	"mod": func(aInt, bInt int) int { return aInt % bInt },
}

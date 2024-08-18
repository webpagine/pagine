// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import (
	"path"
)

type v1Path struct{}

func (v1Path) Join(paths ...string) string { return path.Join(paths...) }

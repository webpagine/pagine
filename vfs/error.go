// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package vfs

import "io/fs"

type PathError = fs.PathError

type IllegalPathError struct {
}

func (e *IllegalPathError) Error() string {
	return "illegal path"
}

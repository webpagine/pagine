// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package common

import (
	. "path"
)

type V1Path struct{}

func (V1Path) Base(path string) string { return Base(path) }

func (V1Path) Clean(path string) string { return Clean(path) }

func (V1Path) Dir(path string) string { return Dir(path) }

func (V1Path) Ext(path string) string { return Ext(path) }

func (V1Path) IsAbs(path string) bool { return IsAbs(path) }

func (V1Path) Join(paths ...string) string { return Join(paths...) }

func (V1Path) Match(pattern, name string) (matched bool, valid bool) {
	matched, err := Match(pattern, name)
	if err != nil {
		return false, false
	}

	return matched, true
}

func (V1Path) Split(path string) (dir, file string) { return Split(path) }

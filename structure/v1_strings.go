// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import "strings"

type v1Strings struct{}

func (v1Strings) HasPrefix(str string, prefix string) bool    { return strings.HasPrefix(str, prefix) }
func (v1Strings) TrimPrefix(str string, prefix string) string { return strings.TrimPrefix(str, prefix) }

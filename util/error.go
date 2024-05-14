// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package util

import "fmt"

type ErrorSet struct {
	Errors []error
}

func (e *ErrorSet) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

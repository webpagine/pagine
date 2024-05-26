// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import "fmt"

type MissingRequiredFieldError struct {
	Field string
}

func (e *MissingRequiredFieldError) Error() string {
	return fmt.Sprint("missing required field: ", e.Field)
}

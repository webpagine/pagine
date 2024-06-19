// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import "fmt"

type TemplateNotFoundError struct {
	Template *Template
	Want     string
}

func (e *TemplateNotFoundError) Error() string {
	return fmt.Sprint("template not found in [", e.Template.CanonicalName, "]: ", e.Want)
}

type TemplateUndefinedError struct {
	Name string
}

func (e *TemplateUndefinedError) Error() string {
	return fmt.Sprint("template undefined: ", e.Name)
}

type RecursiveInvokeError struct {
	Templates []string
}

func (e *RecursiveInvokeError) Error() string {
	return fmt.Sprint("recursive template invoke detected")
}

// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

import "fmt"

type TemplateReport struct{}

type UndefinedStdError struct {
	Std string
}

func (e UndefinedStdError) Error() string {
	return fmt.Sprint("undefined standard func set: ", e.Std)
}

type TemplateNotFoundError struct {
	Template *Template
	Key      string
}

func (e *TemplateNotFoundError) Error() string {
	return fmt.Sprint("template key [", e.Key, "] not found in [", e.Template.CanonicalName, "]")
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

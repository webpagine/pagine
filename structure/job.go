// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package structure

const (
	_ = iota

	WAITING // Waiting for other content being generated.
	FAILED  // Error occurred.
	DONE    // Content generated successfully.
)

type Job struct {
	Requires []string `toml:"requires"`
}

func (j *Job) Execute() error {
	// TODO
	return nil
}

// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package workflow

import (
	"github.com/jellyterra/collection-go"
)

type Job struct {
	Title string

	Commands []*Command
}

func (j *Job) Run() (*JobReport, error) {
	var (
		reports collection.Vector[*CommandReport]
	)

	for _, cmd := range j.Commands {
		reports.Push(&CommandReport{
			Command: cmd,
			Err:     cmd.Execute(),
		})
	}

	return &JobReport{
		Job:            j,
		CommandReports: reports.Raw,
	}, nil
}

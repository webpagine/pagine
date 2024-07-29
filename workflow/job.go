// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package workflow

import (
	"github.com/jellyterra/collection-go"
	"sync"
)

type Job struct {
	Commands []*Command
}

func (j *Job) Run() (JobReport, error) {
	var (
		reports collection.SyncVector[CommandReport]

		wg sync.WaitGroup
	)

	for _, cmd := range j.Commands {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reports.Push(CommandReport{
				Command: cmd,
				Err:     cmd.Execute(),
			})
		}()
	}

	wg.Wait()

	return JobReport{
		Job:            j,
		CommandReports: reports.It.Raw,
	}, nil
}

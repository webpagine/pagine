// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package workflow

import (
	"github.com/jellyterra/collection-go"
	"sync"
)

type Stage struct {
	Jobs []*Job
}

func (s *Stage) Run() (*StageReport, error) {
	var (
		jobReports, jobFailures collection.SyncVector[JobReport]

		wg sync.WaitGroup
	)

	for _, job := range s.Jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()

			report, err := job.Run()
			if err != nil {
				jobFailures.Push(JobReport{
					Job:     job,
					Failure: err,
				})
				return
			}
			jobReports.Push(report)
		}()
	}

	wg.Wait()

	return &StageReport{
		Stage:      s,
		JobReports: jobReports.It.Raw,
	}, nil
}

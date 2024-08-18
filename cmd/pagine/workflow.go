// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/webpagine/pagine/v2/config"
	"github.com/webpagine/pagine/v2/vfs"
)

func RunWorkflow(root *vfs.DirFS) error {
	wf, err := config.LoadWorkflow(root)
	if err != nil {
		return err
	}

	for stageIndex, stage := range wf.Stages {
		report, err := stage.Run()
		if err != nil {
			return err
		}

		for _, jobReport := range report.JobReports {
			for _, cmdReport := range jobReport.CommandReports {
				if cmdReport.Err != nil {
					fmt.Printf("Stage %d - Job failed \"%s\": %e\n", stageIndex, jobReport.Job.Title, cmdReport.Err)
				}
			}
		}
	}

	return nil
}

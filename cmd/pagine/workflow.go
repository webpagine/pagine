// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/config"
	"github.com/webpagine/pagine/v2/vfs"
	"github.com/webpagine/pagine/v2/workflow"
	"io/fs"
	"sync"
)

func CollectAndRunWorkflows(root, dest *vfs.DirFS, jobBuilderRoot fs.StatFS) error {
	bases, err := CollectWorkflows(root)
	if err != nil {
		return err
	}

	if len(bases) == 0 {
		return nil
	}

	var workflows []*workflow.Workflow

	for _, base := range bases {
		origin, err := root.Chroot(base)
		if err != nil {
			return err
		}
		wfRoot, err := dest.Chroot(base)
		if err != nil {
			return err
		}

		wf, err := config.LoadWorkflow(origin, wfRoot, jobBuilderRoot)
		if err != nil {
			return err
		}

		workflows = append(workflows, wf)
	}

	var wg sync.WaitGroup

	for _, wf := range workflows {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := RunWorkflow(wf)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	wg.Wait()

	fmt.Println("Workflows complete.")

	return nil
}

func CollectWorkflows(root *vfs.DirFS) (dirs []string, _ error) {
	return dirs, fs.WalkDir(root, "/", func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		_, err = root.Stat(path + "/workflow.yaml")
		if err != nil {
			return nil // Skip.
		}

		dirs = append(dirs, path)

		return nil
	})
}

func RunWorkflow(wf *workflow.Workflow) error {
	var stageReports collection.Vector[*workflow.StageReport]

	for _, stage := range wf.Stages {
		stageReport, err := stage.Run()
		if err != nil {
			return err
		}

		stageReports.Push(stageReport)
	}

	for stageIndex, stageReport := range stageReports.Raw {
		for _, jobReport := range stageReport.JobReports {
			for _, cmdReport := range jobReport.CommandReports {

				switch {
				case cmdReport.Err != nil:
					fmt.Println("\n===== Stage", stageIndex, "-", jobReport.Job.Title, ":", cmdReport.Err.Error())
					fmt.Println(cmdReport.Command.Exec, cmdReport.Command.Args)
					fmt.Print("\n", cmdReport.Output.String())
				case cmdReport.Output.Len() == 0:
				default:
					fmt.Println("\n===== Stage", stageIndex, "-", jobReport.Job.Title)
					fmt.Print(cmdReport.Output.String())
				}
			}
		}
	}

	return nil
}

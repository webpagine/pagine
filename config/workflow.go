// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package config

import (
	"fmt"
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/vfs"
	"github.com/webpagine/pagine/v2/workflow"
	"github.com/webpagine/pagine/v2/workflow/builtin"
	"io/fs"
	"os"
)

type TscJobConfig struct {
	Path string `json:"path"`
}

type Workflow struct {
	Stage []struct {
		Title string           `yaml:"title"`
		Job   []map[string]any `yaml:"job"`
	} `yaml:"stage"`
}

func LoadJob(root *vfs.DirFS, m map[string]any, jobBuilderRoot fs.StatFS) (*workflow.Job, error) {
	if m["type"] == nil {
		return nil, fmt.Errorf("job type is not specified")
	}
	jobBuilderFile := m["type"].(string) + ".tmpl"

	_, err := jobBuilderRoot.Stat(jobBuilderFile)
	switch {
	case err == nil:
	case os.IsNotExist(err):
		return workflow.BuildJob(builtin.BuiltinBuilders, jobBuilderFile, root, m)
	default:
		return nil, err
	}

	return workflow.BuildJob(jobBuilderRoot, jobBuilderFile, root, m)
}

func LoadWorkflow(root *vfs.DirFS, jobBuilderRoot fs.StatFS) (*workflow.Workflow, error) {

	var (
		rawWorkflow Workflow

		stages collection.Vector[*workflow.Stage]
	)

	err := UnmarshalYAMLFile(root, "workflow.yaml", &rawWorkflow)
	if err != nil {
		return nil, err
	}

	for _, stage := range rawWorkflow.Stage {
		var jobs collection.Vector[*workflow.Job]

		for _, job := range stage.Job {
			job, err := LoadJob(root, job, jobBuilderRoot)
			if err != nil {
				return nil, err
			}

			jobs.Push(job)
		}

		stages.Push(&workflow.Stage{Jobs: jobs.Raw})
	}

	return &workflow.Workflow{
		Stages: stages.Raw,
	}, nil
}

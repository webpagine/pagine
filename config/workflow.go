// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package config

import (
	"fmt"
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/vfs"
	"github.com/webpagine/pagine/v2/workflow"
	"github.com/webpagine/pagine/v2/workflow/builtin/tsc"
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

func LoadJob(root *vfs.DirFS, m map[string]any) (*workflow.Job, error) {
	gen, ok := jobGens[m["type"].(string)]
	if !ok {
		return nil, fmt.Errorf("unknown job type: %s", m["type"])
	}

	job, err := gen(root, m)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func LoadWorkflow(root *vfs.DirFS) (*workflow.Workflow, error) {

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
			job, err := LoadJob(root, job)
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

type TscJob struct {
	Path string `json:"path"`
}

func generateTscJob(root *vfs.DirFS, m map[string]any) (*workflow.Job, error) {
	var tscJob TscJob

	err := UnmarshalMap(m, &tscJob)
	if err != nil {
		return nil, err
	}

	cmd, err := tsc.BuildTS(root, tscJob.Path)

	return &workflow.Job{
		Title:    m["title"].(string),
		Commands: []*workflow.Command{cmd},
	}, nil
}

var jobGens = map[string]func(root *vfs.DirFS, m map[string]any) (*workflow.Job, error){
	"tsc/v1": generateTscJob,
}

// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package workflow

import (
	"bytes"
	"os/exec"
)

type Command struct {
	Exec string   `yaml:"exec"`
	Args []string `yaml:"args"`
}

func (c *Command) Execute() *CommandReport {
	output := bytes.NewBuffer(nil)

	cmd := exec.Command(c.Exec, c.Args...)
	cmd.Stdout, cmd.Stderr = output, output

	return &CommandReport{
		Command: c,
		Err:     cmd.Run(),
		Output:  output,
	}
}

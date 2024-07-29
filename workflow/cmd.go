// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package workflow

import "os/exec"

type Command struct {
	Exec string   `toml:"exec"`
	Args []string `toml:"args"`
}

func (c *Command) Execute() error {
	return exec.Command(c.Exec, c.Args...).Run()
}

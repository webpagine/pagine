// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package tsc

import (
	"github.com/webpagine/pagine/v2/vfs"
	"github.com/webpagine/pagine/v2/workflow"
)

func BuildTS(fs *vfs.DirFS, path string) (*workflow.Command, error) {
	return &workflow.Command{
		Exec: "/bin/tsc",
		Args: []string{"--project", fs.Path + "/tsconfig.json"},
	}, nil
}

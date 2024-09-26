// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package workflow

import (
	"bytes"
	"fmt"
	"github.com/jellyterra/collection-go"
	"github.com/webpagine/pagine/v2/common"
	"github.com/webpagine/pagine/v2/vfs"
	"io/fs"
	"path"
	"text/template"
)

func BuildJob(builderTmplRoot fs.FS, builderTmpl string, origin, root *vfs.DirFS, config map[string]any) (*Job, error) {
	var args common.Slice
	var export = common.Map{
		Raw: map[string]any{
			"title":      "",
			"executable": "",
			"args":       &args,
		},
	}

	funcMap := map[string]any{
		"origin":  common.WrapObject(origin.Path),
		"root":    common.WrapObject(root.Path),
		"job":     common.WrapObject(&export),
		"path":    common.WrapObject(common.V1Path{}),
		"strings": common.WrapObject(common.V1Strings{}),
	}
	collection.MergeRawMap(funcMap, common.V1Arithmetic)
	collection.MergeRawMap(funcMap, common.V1Builtin)

	t, err := template.New(path.Base(builderTmpl)).Funcs(funcMap).ParseFS(builderTmplRoot, builderTmpl)
	if err != nil {
		return nil, err
	}
	err = t.Execute(bytes.NewBuffer(nil), config)
	if err != nil {
		return nil, err
	}

	switch {
	case !export.Has("title"):
		return nil, fmt.Errorf("missing property: title")
	case !export.Has("executable"):
		return nil, fmt.Errorf("missing property: executable")
	}

	var (
		title      = export.Raw["title"].(string)
		executable = export.Raw["executable"].(string)
	)

	return &Job{
		Title: title,
		Commands: []*Command{
			{Exec: executable, Args: args.StringSlice()},
		},
	}, nil
}

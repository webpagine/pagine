// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/webpagine/pagine/v2/structure"
	"github.com/webpagine/pagine/v2/vfs"
	"sync"
)

func ExecuteLevel(env *structure.Env, dest *vfs.DirFS, level *structure.Level) {
	report, err := level.Generate(env, dest)
	if err != nil {
		fmt.Print(`Executing level "`, level.Root.Path, `" failed:`, err, "\n")
	}

	for _, unitReport := range report.UnitReports {
		if unitReport.Error != nil {
			fmt.Print("\t\"", unitReport.Unit.Output, `": `, unitReport.Error, "\n")
			continue
		}

		if unitReport.TemplateErrors != nil {
			fmt.Print(`Unit "`, unitReport.Unit.Output, `" has template errors:`, "\n")
		}

		for _, templateError := range unitReport.TemplateErrors {
			fmt.Println("\t", templateError)
		}
	}
}

func ExecuteLevels(env *structure.Env, dest *vfs.DirFS, levels ...*structure.Level) {
	var wg sync.WaitGroup

	for _, level := range levels {
		wg.Add(1)
		go func() {
			defer wg.Done()

			go ExecuteLevels(env, dest, level.Levels...)
			ExecuteLevel(env, dest, level)
		}()
	}

	wg.Wait()
}

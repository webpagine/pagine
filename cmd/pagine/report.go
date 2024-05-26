// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/webpagine/go-pagine/structure"
)

func PrintReport(report *structure.SiteGenerationReport) {
	if report.FileSystemErrors.RawMap != nil {
		fmt.Println("File system errors:")
		for path, err := range report.FileSystemErrors.RawMap {
			fmt.Print("\t[", path, "]\t", err, "\n")
		}
	}

	if report.PageUnmarshalErrors.RawMap != nil {
		fmt.Println("Page unmarshal errors:")
		for path, err := range report.PageUnmarshalErrors.RawMap {
			fmt.Print("\t[", path, "]\t", err, "\n")
		}
	}

	if report.PageGenerationErrors.RawMap != nil {
		fmt.Println("Page generation errors:")
		for path, err := range report.PageGenerationErrors.RawMap {
			fmt.Print("\t[", path, "]\t", err, "\n")
		}
	}
}

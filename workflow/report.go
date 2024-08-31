// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package workflow

import "bytes"

type CommandReport struct {
	Command *Command

	Err error

	Output *bytes.Buffer
}

type JobReport struct {
	Job *Job

	Failure error

	CommandReports []*CommandReport
}

type StageReport struct {
	Stage *Stage

	JobReports []*JobReport
}

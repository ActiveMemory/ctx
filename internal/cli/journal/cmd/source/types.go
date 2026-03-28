//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

// Opts holds all flags for the source subcommand.
type Opts struct {
	ShowID      string
	Latest      bool
	Full        bool
	Limit       int
	Project     string
	Tool        string
	Since       string
	Until       string
	AllProjects bool
}

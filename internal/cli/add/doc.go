//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add hosts the shared core libraries that power
// every noun-first add subcommand:
//
//   - ctx task add        (internal/cli/task/cmd/add)
//   - ctx decision add    (internal/cli/decision/cmd/add)
//   - ctx learning add    (internal/cli/learning/cmd/add)
//   - ctx convention add  (internal/cli/convention/cmd/add)
//
// The directory contains no Go source at the top level. Its
// child package core/ groups validation, content extraction,
// markdown formatting, section-aware insertion, and section
// normalization, plus the build/ helper that assembles a
// noun-bound cobra command and the run/ entry point that
// executes the add pipeline. The verb-first ctx add parent
// was retired by specs/cli-add-symmetry.md; this directory
// remains as a logical home for the shared add machinery.
package add

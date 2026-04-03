//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package show provides the "ctx trace" CLI subcommand (the default
// trace view).
//
// It displays context refs for a specific commit or the last N commits.
// Supports --last to control how many commits to show and --json for
// machine-readable output.
//
// Key exports: [Cmd], [Run].
// Used by the trace command tree as the root trace display handler.
package show

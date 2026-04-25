//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// GitCommit matches git commit commands.
var GitCommit = regexp.MustCompile(`git\s+commit`)

// GitAmend matches the --amend flag.
var GitAmend = regexp.MustCompile(`--amend`)

// TaskRef matches Phase-style task references like HA.1, P-2.5, PD.3, CT.1.
var TaskRef = regexp.MustCompile(`\b[A-Z]+-?\d+\.?\d*\b`)

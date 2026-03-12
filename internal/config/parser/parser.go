//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

// Parser configuration.
const (
	// LinesToPeek is the number of lines to scan when detecting file format.
	LinesToPeek = 50
	// DirSubagents is the directory name for sidechain sessions that share
	// the parent sessionId and would cause duplicates if scanned.
	DirSubagents = "subagents"
)

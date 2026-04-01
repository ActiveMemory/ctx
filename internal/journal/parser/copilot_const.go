//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

// Copilot JSONL line Kind values. Parser-internal: these are wire
// format discriminators, not configurable.
const (
	// copilotKindSnapshot is a full session snapshot (kind=0).
	copilotKindSnapshot = 0
	// copilotKindScalarPatch is a scalar field replacement (kind=1).
	copilotKindScalarPatch = 1
	// copilotKindObjectPatch is an array/object replacement (kind=2).
	copilotKindObjectPatch = 2
)

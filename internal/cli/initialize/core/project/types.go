//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

import "github.com/ActiveMemory/ctx/internal/config/dir"

// Dirs lists the project-root directories created by ctx init,
// each with an explanatory README.md.
var Dirs = []string{
	dir.Specs,
	dir.Ideas,
}

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/config/dir"
)

// Gitignore lists the recommended .gitignore entries added by ctx init.
var Gitignore = []string{
	path.Join(dir.Context, dir.Journal, "/"),
	path.Join(dir.Context, dir.JournalSite, "/"),
	path.Join(dir.Context, dir.JournalObsidian, "/"),
	path.Join(dir.Context, dir.Logs, "/"),
	".context/.ctx.key",
	".context/.context.key",
	".context/.scratchpad.key",
	".context/state/",
	".claude/settings.local.json",
}

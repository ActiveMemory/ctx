//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	cfgHandover "github.com/ActiveMemory/ctx/internal/config/handover"
)

// FileGitignore is the .gitignore filename.
const FileGitignore = ".gitignore"

// GitignoreHeader is the section comment prepended to ctx-managed entries.
const GitignoreHeader = "# ctx managed entries"

// Gitignore lists the recommended .gitignore entries added by ctx init.
//
// The `handovers/` carve-out ignores per-session artifacts (which
// grow unbounded and carry operator-specific identifiers) while
// keeping `.gitkeep` tracked so the read-side missing-dir gate
// (.context/handovers/ missing → ctx init --upgrade) passes for
// fresh clones.
var Gitignore = []string{
	path.Join(dir.Context, dir.Journal, "/"),
	path.Join(dir.Context, dir.JournalSite, "/"),
	path.Join(dir.Context, dir.JournalObsidian, "/"),
	path.Join(dir.Context, dir.Logs, "/"),
	".context/.ctx.key",
	".context/state/",
	path.Join(dir.Context, cfgHandover.Subdir, "*"),
	"!" + path.Join(dir.Context, cfgHandover.Subdir, ".gitkeep"),
	".claude/settings.local.json",
}

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
	cfgProposal "github.com/ActiveMemory/ctx/internal/config/proposal"
)

// FileGitignore is the .gitignore filename.
const FileGitignore = ".gitignore"

// GitignoreHeader is the section comment prepended to ctx-managed entries.
const GitignoreHeader = "# ctx managed entries"

// Gitignore lists the recommended .gitignore entries added by ctx init.
//
// The `handovers/` and `proposals/` carve-outs ignore per-session
// AI-generated artifacts (which grow unbounded and may carry
// operator-specific identifiers) while keeping `.gitkeep` tracked
// so the read-side missing-dir gates pass for fresh clones. The
// `proposals/` entry mirrors the handover shape; DECISIONS
// 2026-05-22-220000 (Phase BE Task 1) named this carve-out
// alongside the proposal-queue writer in
// `internal/write/proposal/`.
var Gitignore = []string{
	path.Join(dir.Context, dir.Journal, "/"),
	path.Join(dir.Context, dir.JournalSite, "/"),
	path.Join(dir.Context, dir.JournalObsidian, "/"),
	path.Join(dir.Context, dir.Logs, "/"),
	".context/.ctx.key",
	".context/state/",
	path.Join(dir.Context, cfgHandover.Subdir, "*"),
	"!" + path.Join(dir.Context, cfgHandover.Subdir, ".gitkeep"),
	path.Join(dir.Context, cfgProposal.Subdir, "*"),
	"!" + path.Join(dir.Context, cfgProposal.Subdir, ".gitkeep"),
	".claude/settings.local.json",
}

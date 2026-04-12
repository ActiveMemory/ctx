//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for Claude plugin check display labels.
const (
	// DescKeyWriteClaudecheckUnknown is the placeholder for
	// unknown fields.
	DescKeyWriteClaudecheckUnknown = "write.claudecheck-unknown"
	// DescKeyWriteClaudecheckNone is the placeholder for
	// empty fields.
	DescKeyWriteClaudecheckNone = "write.claudecheck-none"
	// DescKeyWriteClaudecheckSourceDir is the label for
	// directory-sourced installs.
	DescKeyWriteClaudecheckSourceDir = "write.claudecheck-source-dir"
	// DescKeyWriteClaudecheckSourceGitHub is the label for
	// marketplace installs.
	DescKeyWriteClaudecheckSourceGitHub = "write.claudecheck-source-github"
	// DescKeyWriteClaudecheckEnabledBoth is the label when
	// enabled in both scopes.
	DescKeyWriteClaudecheckEnabledBoth = "write.claudecheck-enabled-both"
	// DescKeyWriteClaudecheckEnabledGlobal is the label when
	// enabled globally.
	DescKeyWriteClaudecheckEnabledGlobal = "write.claudecheck-enabled-global"
	// DescKeyWriteClaudecheckEnabledLocal is the label when
	// enabled per-project.
	DescKeyWriteClaudecheckEnabledLocal = "write.claudecheck-enabled-local"
	// DescKeyWriteClaudecheckVersionOpen is the opening
	// bracket before the SHA.
	DescKeyWriteClaudecheckVersionOpen = "write.claudecheck-version-open"
	// DescKeyWriteClaudecheckVersionClose is the closing
	// bracket after the SHA.
	DescKeyWriteClaudecheckVersionClose = "write.claudecheck-version-close"
)

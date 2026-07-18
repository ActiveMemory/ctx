//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for progressive-disclosure output labels. Each value is the
// key of an entry in commands/text/write.yaml.
const (
	// DescKeyWriteDisclosureKind labels the inspected root's kind.
	DescKeyWriteDisclosureKind = "write.disclosure-kind"
	// DescKeyWriteDisclosureStagedHeader heads the staged-entry list.
	DescKeyWriteDisclosureStagedHeader = "write.disclosure-staged-header"
	// DescKeyWriteDisclosureStagedLine formats one staged entry.
	DescKeyWriteDisclosureStagedLine = "write.disclosure-staged-line"
	// DescKeyWriteDisclosureThemesHeader heads the themes list.
	DescKeyWriteDisclosureThemesHeader = "write.disclosure-themes-header"
	// DescKeyWriteDisclosureThemeLine formats one theme.
	DescKeyWriteDisclosureThemeLine = "write.disclosure-theme-line"
	// DescKeyWriteDisclosureNone marks an empty list.
	DescKeyWriteDisclosureNone = "write.disclosure-none"
)

// DescKeys for progressive-disclosure guard and invariant errors.
// Each value is the key of an entry in commands/text/errors.yaml; the
// DescKey <-> YAML mapping is a bijection enforced by
// internal/audit.TestDescKeyYAMLLinkage.
const (
	// DescKeyErrDisclosureBrokenThemeLink: a theme link in the root
	// resolves to a nonexistent path.
	DescKeyErrDisclosureBrokenThemeLink = "err.disclosure.broken-theme-link"
	// DescKeyErrDisclosureDuplicateEntry: an entry appears in more than
	// one place (staging and/or multiple theme files).
	DescKeyErrDisclosureDuplicateEntry = "err.disclosure.duplicate-entry"
	// DescKeyErrDisclosureEntryBelowThemes: a "## [" entry sits below the
	// "## Themes" section instead of in the staging zone above it.
	DescKeyErrDisclosureEntryBelowThemes = "err.disclosure.entry-below-themes"
	// DescKeyErrDisclosureMissingThemeFile: a theme gist in the root has
	// no matching theme file.
	DescKeyErrDisclosureMissingThemeFile = "err.disclosure.missing-theme-file"
	// DescKeyErrDisclosureMultipleThemes: the root has more than one
	// "## Themes" section.
	DescKeyErrDisclosureMultipleThemes = "err.disclosure.multiple-themes"
	// DescKeyErrDisclosureNotKnowledgeFileMsg: the sentinel text for a
	// file that is not a canonical knowledge file.
	DescKeyErrDisclosureNotKnowledgeFileMsg = "err.disclosure.not-knowledge-file-msg"
	// DescKeyErrDisclosureNotKnowledgeFile: the format wrapper naming the
	// offending path for a not-a-knowledge-file error.
	DescKeyErrDisclosureNotKnowledgeFile = "err.disclosure.not-knowledge-file"
	// DescKeyErrDisclosureOrphanThemeFile: a theme file has no matching
	// gist in the root.
	DescKeyErrDisclosureOrphanThemeFile = "err.disclosure.orphan-theme-file"
	// DescKeyErrDisclosureStagingUnparsable: the staging zone could not
	// be parsed into discrete entries.
	DescKeyErrDisclosureStagingUnparsable = "err.disclosure.staging-unparsable"
)

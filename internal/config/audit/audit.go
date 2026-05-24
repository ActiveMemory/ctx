//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"time"

	"github.com/ActiveMemory/ctx/internal/config/file"
)

// Audit-channel filesystem layout (under `.context/`).
const (
	// DirName is the per-project audit directory: a sibling
	// of the canonical context files inside `.context/`.
	DirName = "audit"

	// DismissedLedger is the basename of the dismissal ledger
	// (kept under DirName so that nuking `.context/state/`
	// does not silently re-surface dismissed audits).
	DismissedLedger = ".dismissed.json"

	// ReportExt is the on-disk extension for an audit report
	// (`<kind><ReportExt>`). Markdown so a human can open and
	// read the body directly without ceremony.
	ReportExt = file.ExtMarkdown
)

// Frontmatter and body shape constants.
const (
	// FrontmatterDelimiter is the YAML frontmatter fence
	// used by audit report files.
	FrontmatterDelimiter = "---"

	// StatusFindings indicates the audit produced
	// actionable findings the hook should relay.
	StatusFindings = "findings"

	// StatusClean indicates the audit found no issues; the
	// hook stays silent for `clean` reports.
	StatusClean = "clean"
)

// Staleness thresholds.
const (
	// StalenessAge is the wall-clock age beyond which a
	// report is considered too old to relay without a
	// `STALE` prefix. Matches the typical user feedback
	// loop on a manually-triggered out-of-band auditor.
	StalenessAge = 7 * 24 * time.Hour
)

// Template variable keys for the check-audit hook.
const (
	// VarList is the template variable for the formatted
	// list of audit reports rendered inside the relay box.
	VarList = "AuditList"
)

// Render-format constants for the check-audit hook body
// and the `ctx audit list` output.
const (
	// FmtReportHeader formats one report's id + commit-range
	// header line. Takes (id, commit-range, newline).
	FmtReportHeader = "[%s] %s%s"

	// FmtStaleLine wraps the STALE prefix line with a
	// trailing newline. Takes (prefix-text, newline).
	FmtStaleLine = "%s%s"

	// FmtAgeDays formats an age in whole days for the
	// STALE prefix (e.g. "9d"). Takes (days).
	FmtAgeDays = "%dd"

	// FmtAgeHours formats an age in whole hours for the
	// STALE prefix when the report is younger than a
	// day. Takes (hours).
	FmtAgeHours = "%dh"

	// SuffixDismissed is appended to a report's status
	// column in `ctx audit list` when the report is
	// dismissed against its current digest.
	SuffixDismissed = " (dismissed)"
)

// HoursPerDay is the divisor used when humanizing a
// report's age into whole days. Pulled out as a named
// constant so the magic-values audit gate stays clean.
const HoursPerDay = 24

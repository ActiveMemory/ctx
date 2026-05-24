//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for the check-audit hook (verbatim relay of
// out-of-band audit reports).
const (
	// DescKeyCheckAuditBoxTitle is the text key for the
	// check-audit nudge-box title.
	DescKeyCheckAuditBoxTitle = "check-audit.box-title"
	// DescKeyCheckAuditDismissHint is the text key for the
	// per-id dismissal hint surfaced in the relay box.
	DescKeyCheckAuditDismissHint = "check-audit.dismiss-hint"
	// DescKeyCheckAuditDismissAllHint is the text key for the
	// dismiss-all hint surfaced in the relay box.
	DescKeyCheckAuditDismissAllHint = "check-audit.dismiss-all-hint"
	// DescKeyCheckAuditNudgeFormat is the text key for the
	// human-readable relay-message suffix
	// (e.g. "You have N unread audit reports...").
	DescKeyCheckAuditNudgeFormat = "check-audit.nudge-format"
	// DescKeyCheckAuditRelayPrefix is the text key for the
	// verbatim-relay prefix string.
	DescKeyCheckAuditRelayPrefix = "check-audit.relay-prefix"
	// DescKeyCheckAuditReportSeparator is the text key for the
	// rule drawn between multiple reports inside one relay box.
	DescKeyCheckAuditReportSeparator = "check-audit.report-separator"
	// DescKeyCheckAuditStalePrefix is the text key for the
	// per-report STALE prefix. Takes (audited-at, branch-tip).
	DescKeyCheckAuditStalePrefix = "check-audit.stale-prefix"
)

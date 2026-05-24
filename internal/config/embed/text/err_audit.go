//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for audit-channel errors.
const (
	// DescKeyErrAuditReadReport is the text key for an audit
	// report read failure.
	DescKeyErrAuditReadReport = "err.audit.read-report"
	// DescKeyErrAuditParseReport is the text key for an audit
	// report frontmatter parse failure.
	DescKeyErrAuditParseReport = "err.audit.parse-report"
	// DescKeyErrAuditWriteDismissal is the text key for a
	// dismissal ledger write failure.
	DescKeyErrAuditWriteDismissal = "err.audit.write-dismissal"
	// DescKeyErrAuditReadDismissal is the text key for a
	// dismissal ledger read failure.
	DescKeyErrAuditReadDismissal = "err.audit.read-dismissal"
	// DescKeyErrAuditUnknownID is the text key for an unknown
	// audit id at the CLI.
	DescKeyErrAuditUnknownID = "err.audit.unknown-id"
	// DescKeyErrAuditIDRequired is the text key for a missing
	// audit id at the dismiss CLI.
	DescKeyErrAuditIDRequired = "err.audit.id-required"
	// DescKeyErrAuditNoFrontmatter is the text key for the
	// sentinel "missing YAML frontmatter" parse error.
	DescKeyErrAuditNoFrontmatter = "err.audit.no-frontmatter"
	// DescKeyErrAuditUnterminatedFrontmatter is the text key
	// for the sentinel "unterminated frontmatter" parse
	// error.
	DescKeyErrAuditUnterminatedFrontmatter = "err.audit.unterminated-frontmatter"
)

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for post-commit hooks.
const (
	// DescKeyPostCommitFallback is the text key for post commit fallback messages.
	DescKeyPostCommitFallback = "post-commit.fallback"
	// DescKeyPostCommitRelayMessage is the text key for post commit relay message
	// messages.
	DescKeyPostCommitRelayMessage = "post-commit.relay-message"
	// DescKeyPostCommitRelayPrefix is the text key for post commit relay prefix
	// messages.
	DescKeyPostCommitRelayPrefix = "post-commit.relay-prefix"
	// DescKeyPostCommitMissingSpec is the text key for post commit missing spec
	// messages.
	DescKeyPostCommitMissingSpec = "post-commit.missing-spec"
	// DescKeyPostCommitMissingSignoff is the text key for post commit missing
	// signoff messages.
	DescKeyPostCommitMissingSignoff = "post-commit.missing-signoff"
	// DescKeyPostCommitMissingBody is the text key for post commit missing body
	// messages.
	DescKeyPostCommitMissingBody = "post-commit.missing-body"
	// DescKeyPostCommitMissingTaskRef is the text key for post commit missing
	// task ref messages.
	DescKeyPostCommitMissingTaskRef = "post-commit.missing-task-ref"
	// DescKeyPostCommitMissingTaskUpdate is the text key for post commit missing
	// task update messages.
	DescKeyPostCommitMissingTaskUpdate = "post-commit.missing-task-update"
	// DescKeyPostCommitSeverityInformal is the text key for post commit severity
	// informal messages.
	DescKeyPostCommitSeverityInformal = "post-commit.severity-informal"
	// DescKeyPostCommitSeveritySkipped is the text key for post commit severity
	// skipped messages.
	DescKeyPostCommitSeveritySkipped = "post-commit.severity-bypassed"
	// DescKeyPostCommitAuditTitle is the text key for post commit audit title
	// messages.
	DescKeyPostCommitAuditTitle = "post-commit.audit-title"
	// DescKeyPostCommitAuditContent is the text key for post commit audit content
	// messages.
	DescKeyPostCommitAuditContent = "post-commit.audit-content"
)

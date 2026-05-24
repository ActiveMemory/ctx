//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for audit-channel writer output.
const (
	// DescKeyWriteAuditNone is the text key for the
	// "no audit reports" message.
	DescKeyWriteAuditNone = "write.audit-none"
	// DescKeyWriteAuditListItem is the text key for the
	// per-report row in `ctx audit list` output.
	// Takes (id, status, commit-range, generated-at).
	DescKeyWriteAuditListItem = "write.audit-list-item"
	// DescKeyWriteAuditDismissed is the text key for the
	// per-id dismissal confirmation.
	DescKeyWriteAuditDismissed = "write.audit-dismissed"
	// DescKeyWriteAuditDismissedAll is the text key for the
	// count-of-N dismissal confirmation under --all.
	DescKeyWriteAuditDismissedAll = "write.audit-dismissed-all"
)

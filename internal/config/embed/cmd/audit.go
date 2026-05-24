//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for the audit command and its subcommands.
const (
	// UseAudit is the cobra Use string for the audit command.
	UseAudit = "audit"
	// UseAuditList is the cobra Use string for `audit list`.
	UseAuditList = "list"
	// UseAuditShow is the cobra Use string for `audit show`.
	UseAuditShow = "show ID"
	// UseAuditDismiss is the cobra Use string for `audit
	// dismiss`. Accepts one or more IDs, or --all.
	UseAuditDismiss = "dismiss [ID...]"
)

// DescKeys for audit-channel commands.
const (
	// DescKeyAudit is the description key for the audit
	// command.
	DescKeyAudit = "audit"
	// DescKeyAuditList is the description key for `audit
	// list`.
	DescKeyAuditList = "audit.list"
	// DescKeyAuditShow is the description key for `audit
	// show`.
	DescKeyAuditShow = "audit.show"
	// DescKeyAuditDismiss is the description key for `audit
	// dismiss`.
	DescKeyAuditDismiss = "audit.dismiss"
)

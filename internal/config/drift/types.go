//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

// IssueType categorizes a drift issue for grouping
// and filtering.
type IssueType = string

// Drift issue type constants for categorization.
const (
	// IssueDeadPath indicates a file path reference
	// that no longer exists.
	IssueDeadPath IssueType = "dead_path"
	// IssueStaleness indicates accumulated completed
	// tasks needing archival.
	IssueStaleness IssueType = "staleness"
	// IssueSecret indicates a file that may contain
	// secrets or credentials.
	IssueSecret IssueType = "potential_secret"
	// IssueMissing indicates a required context file
	// that does not exist.
	IssueMissing IssueType = "missing_file"
	// IssueStaleAge indicates a context file that
	// hasn't been modified recently.
	IssueStaleAge IssueType = "stale_age"
	// IssueEntryCount indicates a knowledge file has
	// too many entries.
	IssueEntryCount IssueType = "entry_count"
	// IssueMissingPackage indicates an internal package
	// not documented in ARCHITECTURE.md.
	IssueMissingPackage IssueType = "missing_package"
	// IssueStaleHeader indicates a context file whose
	// comment header doesn't match the embedded template.
	IssueStaleHeader IssueType = "stale_header"
	// IssueInvalidTool indicates an unsupported tool
	// identifier in a steering file or .ctxrc config.
	IssueInvalidTool IssueType = "invalid_tool"
	// IssueHookNoExec indicates a hook script missing
	// the executable permission bit.
	IssueHookNoExec IssueType = "hook_no_exec"
	// IssueStaleSyncFile indicates a synced tool-native
	// file that is out of date compared to its source.
	IssueStaleSyncFile IssueType = "stale_sync_file"
)

// StatusType represents the overall status of a drift
// report.
type StatusType = string

// Drift report status constants.
const (
	// StatusOk means no drift was detected.
	StatusOk StatusType = "ok"
	// StatusWarning means non-critical issues were found.
	StatusWarning StatusType = "warning"
	// StatusViolation means constitution violations were
	// found.
	StatusViolation StatusType = "violation"
)

// CheckName identifies a drift detection check.
type CheckName = string

// Drift detection check name constants.
const (
	// CheckPathReferences validates that file paths in
	// context files exist.
	CheckPathReferences CheckName = "path_references"
	// CheckStaleness detects accumulated completed tasks.
	CheckStaleness CheckName = "staleness_check"
	// CheckConstitution verifies constitution rules are
	// respected.
	CheckConstitution CheckName = "constitution_check"
	// CheckRequiredFiles ensures all required context files
	// are present.
	CheckRequiredFiles CheckName = "required_files"
	// CheckFileAge checks whether context files have been
	// modified recently.
	CheckFileAge CheckName = "file_age_check"
	// CheckEntryCount checks whether knowledge files have
	// excessive entries.
	CheckEntryCount CheckName = "entry_count_check"
	// CheckMissingPackages checks for undocumented internal
	// packages.
	CheckMissingPackages CheckName = "missing_packages"
	// CheckTemplateHeaders checks context file comment
	// headers against templates.
	CheckTemplateHeaders CheckName = "template_headers"
	// CheckSteeringTools validates tool identifiers in
	// steering files.
	CheckSteeringTools CheckName = "steering_tools"
	// CheckHookPerms checks hook scripts for executable
	// permission bits.
	CheckHookPerms CheckName = "hook_permissions"
	// CheckSyncStaleness compares synced tool-native files
	// against source steering files.
	CheckSyncStaleness CheckName = "sync_staleness"
	// CheckRCTool validates the .ctxrc tool field against
	// supported identifiers.
	CheckRCTool CheckName = "rc_tool_field"
)

// Constitution rule names referenced in drift violations.
const (
	// RuleNoSecrets is the constitution rule for secret
	// file detection.
	RuleNoSecrets = "no_secrets"
)

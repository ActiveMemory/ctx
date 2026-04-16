//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift centralizes issue types, status codes,
// check names, and constitution rule identifiers for the
// ctx drift detection subsystem.
//
// Drift detection scans a project's .context/ directory
// for staleness, broken references, missing files, and
// rule violations. Each finding is tagged with an
// IssueType so consumers can filter, group, and render
// results consistently.
//
// # Issue Types
//
// IssueType constants classify each detected problem:
//
//   - IssueDeadPath: a file path reference that no
//     longer resolves on disk
//   - IssueStaleness: completed tasks that should be
//     archived
//   - IssueSecret: a file that may contain secrets
//     or credentials
//   - IssueMissing: a required context file that does
//     not exist
//   - IssueStaleAge: a context file not modified for
//     an extended period
//   - IssueEntryCount: a knowledge file with too many
//     entries
//   - IssueMissingPackage: an internal Go package not
//     documented in ARCHITECTURE.md
//   - IssueStaleHeader: a file whose comment header
//     doesn't match the embedded template
//   - IssueInvalidTool: an unsupported tool identifier
//     in a steering file or .ctxrc
//   - IssueHookNoExec: a hook script missing the
//     executable permission bit
//   - IssueStaleSyncFile: a synced tool-native file
//     that is out of date versus its source
//
// # Status Types
//
// A drift report carries an overall StatusType:
//
//   - StatusOk: no issues detected
//   - StatusWarning: non-critical issues found
//   - StatusViolation: constitution violations found
//
// # Check Names
//
// CheckName constants identify each drift detection
// pass: CheckPathReferences, CheckStaleness,
// CheckConstitution, CheckRequiredFiles, CheckFileAge,
// CheckEntryCount, CheckMissingPackages,
// CheckTemplateHeaders, CheckSteeringTools,
// CheckHookPerms, CheckSyncStaleness, and CheckRCTool.
//
// # Constitution Rules
//
// RuleNoSecrets is the constitution rule name referenced
// when a potential secret file triggers a violation.
//
// # Why Centralized
//
// These constants are shared between the drift scanner,
// the doctor subsystem (which delegates to drift), CLI
// rendering, and fix/auto-repair routines. Centralizing
// them here avoids import cycles and guarantees that
// issue types and check names stay consistent across
// producers and consumers.
package drift

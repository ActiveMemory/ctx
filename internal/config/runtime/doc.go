//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package runtime defines sensible-default numeric
// thresholds that govern ctx's runtime behavior.
//
// Every value here is overridable via the project's
// .ctxrc file, but these constants provide the
// out-of-the-box experience. They control how
// aggressively context is assembled, when nudges
// fire, and when files are considered stale.
//
// # Key Constants
//
//   - [DefaultTokenBudget] (8 000): tokens allocated
//     to the agent context packet.
//   - [DefaultContextWindow] (200 000): assumed
//     model context window size.
//   - [DefaultArchiveAfterDays] (7): days before
//     completed tasks are auto-archived.
//   - [DefaultEntryCountLearnings] (30) and
//     [DefaultEntryCountDecisions] (20): entry
//     count thresholds that trigger consolidation
//     nudges.
//   - [DefaultConventionLineCount] (200): line
//     limit for CONVENTIONS.md before a nudge.
//   - [DefaultInjectionTokenWarn] (15 000): token
//     count that triggers an oversize injection
//     warning.
//   - [DefaultTaskNudgeInterval] (5): Edit/Write
//     calls between task-completion nudges.
//   - [DefaultKeyRotationDays] (90): days before
//     an encryption key rotation nudge.
//   - [DefaultStaleAgeDays] (30): days before a
//     context file is flagged as stale.
//   - [DefaultPruneDays] (7): age threshold for
//     state file pruning.
//
// # Why Centralized
//
// The defaults are consumed by the config loader,
// the drift detector, the agent assembler, and the
// nudge engine. Keeping them in one file makes
// tuning easy and prevents silent disagreement
// between subsystems.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package runtime

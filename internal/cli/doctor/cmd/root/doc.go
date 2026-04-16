//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx doctor" command that
// runs health diagnostics on the project context setup.
//
// # What It Does
//
// Executes a suite of health checks and presents the
// results as a checklist or JSON report. The checks
// cover the full lifecycle of a ctx installation:
//
//   - Context initialization status
//   - Required files presence (TASKS.md, etc.)
//   - .ctxrc validation (syntax, schema)
//   - Context drift detection
//   - Plugin enablement in Claude settings
//   - Companion configuration
//   - Event logging health
//   - Webhook configuration
//   - Reminder validity
//   - Task completion ratios
//   - Context token size vs budget
//   - System resource availability
//   - Recent event activity
//
// # Flags
//
//   - --json, -j: Output results as machine-readable
//     JSON instead of a human-readable checklist.
//
// # Output
//
// In human mode, prints a checklist with pass/warn/
// error icons per check, followed by warning and
// error counts. In JSON mode, outputs a structured
// report with all check results, statuses, and
// summary counts.
//
// # Delegation
//
// [Cmd] builds the cobra.Command with AnnotationSkipInit
// so it runs even when context is not initialized.
// [Run] executes each check function from [core/check],
// tallies warnings and errors, then delegates to
// [core/output] for rendering.
package root

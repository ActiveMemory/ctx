//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core is the umbrella for the doctor command's
// health-check logic.
//
// # Overview
//
// The doctor command runs a series of diagnostic checks
// against the user's context directory, configuration,
// plugins, hooks, and system resources. This package
// groups the check and output sub-packages that contain
// the actual implementations.
//
// # Sub-packages
//
//   - check: individual diagnostic functions, each
//     appending results to a shared [check.Report].
//   - output: renders the report as human-readable
//     text or machine-readable JSON.
//
// # Check Categories
//
// Checks are grouped by category:
//
//   - Structure: context init, required files, ctxrc
//   - Quality: drift detection
//   - Plugin: companion config, plugin enablement
//   - Hooks: event logging, webhook
//   - State: reminders, task completion
//   - Size: context token budget
//   - Events: recent event log activity
//   - Resources: memory, swap, disk, CPU load
//
// The cmd layer calls each Check* function in sequence,
// tallies warnings and errors, then delegates to
// output.Human or output.JSON for rendering.
package core

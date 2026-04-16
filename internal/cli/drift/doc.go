//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift implements the "ctx drift" command for
// detecting stale or invalid context.
//
// The drift command performs a suite of health checks
// against the .context/ directory and reports problems
// that may cause AI agents to operate on outdated or
// broken information. Results can be output as formatted
// text or JSON for scripting.
//
// # Checks Performed
//
//   - Broken path references: file paths mentioned in
//     context files that no longer exist on disk
//   - Staleness indicators: entries whose timestamps
//     exceed configured age thresholds
//   - Constitution violations: rules in CONSTITUTION.md
//     that conflict with other context files
//   - Missing required files: context files expected by
//     the project template but absent from .context/
//
// # Subpackages
//
//   - cmd/root: cobra command definition and flag binding
//   - core: check implementations, path resolution, and
//     staleness evaluation
package drift

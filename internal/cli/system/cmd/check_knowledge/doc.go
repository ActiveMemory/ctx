//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_knowledge implements the
// **`ctx system check-knowledge`** hidden hook, which
// warns when context knowledge files grow beyond their
// configured size thresholds.
//
// # What It Does
//
// The hook checks three context files against their
// respective limits:
//
//   - **DECISIONS.md**: entry count threshold
//   - **LEARNINGS.md**: entry count threshold
//   - **CONVENTIONS.md**: line count threshold
//
// When any file exceeds its configured maximum, the
// hook emits a nudge box recommending consolidation
// (merging duplicates, pruning stale entries). This
// prevents knowledge files from growing so large that
// they become noisy or exceed context window budgets.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On threshold exceeded: a formatted nudge box listing
// which files are over their limits. On all files
// within limits or throttled: no output.
//
// # Throttling
//
// The hook is throttled to fire at most once per day
// using a marker file in the state directory.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble] and delegates the
// health check to [core/knowledge.CheckHealth], which
// handles file scanning, threshold comparison, and
// nudge box formatting in a single call.
package check_knowledge

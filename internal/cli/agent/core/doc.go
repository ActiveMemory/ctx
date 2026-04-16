//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains business logic for the agent
// command, which assembles an AI-optimized context packet
// from project files.
//
// This package is an umbrella that coordinates six
// subpackages, each handling one aspect of the context
// assembly pipeline:
//
// # Budget Allocation (budget/)
//
// The budget subpackage parses the --budget flag, splits
// the token budget across context sections using priority
// weights, and renders the final packet with truncation
// when content exceeds the allocation.
//
// # Cooldown Management (cooldown/)
//
// The cooldown subpackage prevents redundant emissions in
// rapid tool loops by maintaining per-session tombstone
// files with a configurable time-to-live.
//
// # Content Extraction (extract/)
//
// The extract subpackage pulls structured items from
// context files: bullet items, checkbox tasks (checked
// and unchecked), constitution rules, and active tasks
// from TASKS.md.
//
// # Hub Content (hub/)
//
// The hub subpackage loads shared knowledge from the
// .context/hub/ directory, where files received from a
// ctx Hub instance are stored.
//
// # File Ordering (sort/)
//
// The sort subpackage determines the read order for
// context files based on the priority sequence defined
// in config, filtering out empty files.
//
// # Steering and Skills (steering/)
//
// The steering subpackage loads steering files filtered
// by the current tool and resolves named skills from
// the .context/skills/ directory.
//
// # Score (score/)
//
// The score subpackage evaluates context completeness by
// checking which context files exist and are populated,
// producing a numeric health score.
//
// # Data Flow
//
// The cmd/ layer calls into these subpackages to build
// the context packet: sort determines file order, extract
// pulls items, budget allocates space, steering adds
// instructions, and the result is rendered to stdout.
// Cooldown gates the entire pipeline.
package core

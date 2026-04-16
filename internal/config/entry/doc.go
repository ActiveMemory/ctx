//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry defines the vocabulary of context entry
// types, priority levels, and file-mapping helpers used
// throughout the ctx CLI.
//
// Every piece of project context (a task, decision,
// learning, or convention) is identified by a canonical
// type string defined here. The add, complete, and hub
// commands use these strings in switch statements to
// route user input to the correct handler and target
// file.
//
// # Entry Types
//
// The canonical type constants are:
//
//   - Task: a work item stored in TASKS.md
//   - Decision: an architectural decision stored in
//     DECISIONS.md
//   - Learning: a lesson learned stored in
//     LEARNINGS.md
//   - Convention: a code pattern stored in
//     CONVENTIONS.md
//   - Complete: a pseudo-type that marks a task done
//   - Unknown: returned when input doesn't match
//
// AllowedTypes is a set of the four real entry types
// accepted by the hub and add commands.
//
// # Priority Levels
//
// Tasks carry a priority: PriorityHigh, PriorityMedium,
// or PriorityLow. The Priorities slice lists all valid
// levels for shell completion.
//
// # Input Normalization
//
// [FromUserInput] accepts singular and plural forms
// (case-insensitive) and returns the canonical constant.
// For example, "Tasks" and "task" both return Task.
//
// # File Mapping
//
// [CtxFile] and [MustCtxFile] resolve an entry type to
// its target filename (e.g. Decision -> "DECISIONS.md").
// The mapping is thread-safe and backed by config/ctx
// filename constants. MustCtxFile panics on unknown
// types and should only be called after validation.
//
// # Spec Nudge
//
// DefaultSpecSignalWords lists terms in task
// descriptions that suggest a design spec would help.
// SpecNudgeMinLen (150 chars) sets the length threshold
// above which a spec nudge fires regardless of signal
// words. Both are user-configurable via .ctxrc.
//
// # Why Centralized
//
// Entry types appear in CLI commands, hub validation,
// hook payloads, and file I/O routines. Defining them
// in one config package prevents import cycles and
// ensures every consumer agrees on the canonical
// strings.
package entry

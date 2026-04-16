//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package state provides access to the project-scoped
// runtime state directory (.context/state/). Hooks and
// core packages use this directory to persist counters,
// flags, and other ephemeral data that lives alongside
// the context files but is separate from the user-visible
// context.
//
// # Directory Resolution
//
// [Dir] returns the absolute path to .context/state/,
// creating the directory if it does not exist. The
// MkdirAll call is a no-op when the directory is already
// present, so callers do not need to check beforehand.
//
// # Testing Support
//
// [SetDirForTest] overrides the directory returned by
// Dir, redirecting all state reads and writes to a temp
// directory during tests. Pass an empty string to restore
// the default behavior.
//
// # Initialization Guard
//
// [Initialized] reports whether the context directory
// has been properly set up via "ctx init". Hooks check
// this before writing any state to avoid creating a
// partial .context/state/ directory in an uninitialized
// project. When Initialized returns false, hooks should
// silently no-op.
package state

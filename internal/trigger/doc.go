//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger implements ctx's lifecycle
// automation layer: project-owned shell scripts that
// run when an AI session crosses a defined boundary.
//
// # Trigger Types
//
// Six lifecycle events are supported:
//
//   - session-start -- a new AI session begins.
//   - session-end -- an AI session ends.
//   - pre-tool-use -- before a tool call; can block
//     the call via cancel:true.
//   - post-tool-use -- after a tool call completes.
//   - file-save -- a file is saved.
//   - context-add -- a new entry was added to
//     .context/.
//
// Each script receives a JSON [HookInput] on stdin
// and emits a JSON [HookOutput] on stdout.
//
// # Discovery
//
// [Discover] scans .context/hooks/<type>/ and returns
// one [HookInfo] per script, sorted alphabetically.
// The executable permission bit controls whether a
// hook is enabled. [FindByName] locates a single
// script by its stem for enable/disable operations.
//
// # Security
//
// Triggers run with the same privileges as the AI
// tool. The package enforces a strict workflow:
//
//  1. ctx trigger add creates scripts without the
//     executable bit (inert until reviewed).
//  2. ctx trigger enable sets the bit after
//     [ValidatePath] passes.
//  3. [ValidatePath] rejects symlinks, paths that
//     escape the hooks directory, and files lacking
//     the executable bit.
//
// # Execution
//
// [RunAll] runs every enabled hook for a given type
// in alphabetical order. Per-hook behavior:
//
//   - cancel:true halts the chain immediately.
//   - Non-empty context is appended to the aggregate.
//   - Non-zero exit is logged and recorded but does
//     not abort the chain.
//   - Timeout exceeded kills the process group.
//
// The default timeout is 10 seconds
// ([DefaultTimeout]).
//
// # Concurrency
//
// No mutable global state. [RunAll] runs hooks
// sequentially within a single invocation; concurrent
// invocations from different goroutines are safe.
package trigger
